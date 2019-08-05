package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vds/RestaurantManagement/pkg/database"
	"github.com/vds/RestaurantManagement/pkg/encryption"
	"github.com/vds/RestaurantManagement/pkg/middleware"
	"github.com/vds/RestaurantManagement/pkg/models"
	"log"
)
// constants for database queries
const SuperAdminTable="super_admins"
const AdminTable="admins"
const OwnerTable="owners"
const InsertUser="insert into %s(email_id,name,password) values(?,?,?)"
const InsertOwner="insert into owners(email_id,name,password,creator_id) values(?,?,?,?)"
const GetUserIDPassword="select email_id,password from %s where email_id=?"
const DeleteOwner="delete from owners where email_id=?"
const DeleteOwnerByCreator="delete from owners where email_id=? and creator_id=?"



type MySqlDB struct{
	*sql.DB
}
func NewMySqlDB()(*MySqlDB,error){
	db,err:=sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/restaurant")
	if err!=nil{
		return nil,err
	}
	mySqlDB:=&MySqlDB{db}
	return mySqlDB,err
}

//show nearby restaurants
func(db *MySqlDB)ShowNearBy(location *models.Location)(string,error){
	var result string
	rows,err:=db.Query("select JSON_ARRAYAGG(JSON_OBJECT('name',name)) from restaurants where ST_Distance_Sphere(point(lat,lng),point(?,?))/1000 < 10",location.Lat,location.Lng)
	if err!=nil{
		log.Printf("%v",err)
		return "",database.ErrInternal
	}
	rows.Next()
	rows.Scan(&result)
	return result,nil
}



//Interface functions
func (db *MySqlDB)CreateUser(userType string, user *models.User) error{
	var tableName string
	switch userType{
	case middleware.Admin:
		tableName=AdminTable
	case middleware.SuperAdmin:
		tableName=SuperAdminTable
	}
	stmt,err:=db.Prepare(fmt.Sprintf(InsertUser,tableName))
	if err!=nil{
		fmt.Printf("%v",err)
		return database.ErrInternal
	}
	pass,err:=encryption.GenerateHash(user.Password)
	if err!=nil{
		fmt.Printf("%v",err)
		return database.ErrInternal
	}
	_,err=stmt.Exec(user.Email,user.Name,pass)
	if err!=nil{
		fmt.Printf("%v",err)
		return database.ErrDupEmail
	}
	return nil
}

func(db *MySqlDB) LogInUser(userType string,cred *models.Credentials)(string,error){
	var tableName string
	switch userType{
	case middleware.Admin:
		tableName=AdminTable
	case middleware.SuperAdmin:
		tableName=SuperAdminTable
	case middleware.Owner:
		tableName=OwnerTable
	}
	var credOut models.Credentials
	rows,err:=db.Query(fmt.Sprintf(GetUserIDPassword,tableName),cred.Email)
	if err!=nil{
		log.Printf("%v",err)
		return credOut.Email,database.ErrInvalidCredentials
	}
	defer rows.Close()

	rows.Next()
	err=rows.Scan(&credOut.Email,&credOut.Password)
	if err!=nil{
		log.Printf("%v",err)
		return credOut.Email,database.ErrInvalidCredentials
	}
	isValid:=encryption.ComparePasswords(credOut.Password,cred.Password)
	if !isValid{
		return credOut.Email,database.ErrInvalidCredentials
	}
	return credOut.Email,nil
}
//Owner related functions
func(db *MySqlDB)ShowOwners(creatorID string)(string,error){
	var result string
	rows,err:=db.Query("select JSON_ARRAYAGG(JSON_OBJECT('email',email_id,'name', name)) from owners where creator_id=?",creatorID)
	if err!=nil{
		log.Printf("%v",err)
		return "",database.ErrInternal
	}
	rows.Next()
	rows.Scan(&result)
	return result,nil
}

func(db *MySqlDB)ShowOwnersForSuper()(string,error){
	var result string
	rows,err:=db.Query("select JSON_ARRAYAGG(JSON_OBJECT('email',email_id,'name', name)) from owners")
	if err!=nil{
		log.Printf("%v",err)
		return "",database.ErrInternal
	}
	rows.Next()
	rows.Scan(&result)
	return result,nil
}

func(db *MySqlDB)CreateOwners(creatorID string,owners []models.User)error{
	stmt,err:=db.Prepare(InsertOwner)
	if err!=nil{
		log.Printf("%v",err)
		return database.ErrInternal
	}
	for i,owner:=range owners{
		pass,err:=encryption.GenerateHash(owner.Password)
		if err!=nil{
			fmt.Printf("%v",err)
			return database.ErrInternal
		}
		_,err=stmt.Exec(owner.Email,owner.Name,pass,creatorID)
		if err!=nil{
			log.Printf("%v",err)
			go deleteRecordedOwners(db,owners[:i])
			return errors.New(fmt.Sprintf("duplicate entry for ownerID %v" +
				" in entry no. %v ",owner.Email,i+1))
		}
	}
	return nil
}

func(db *MySqlDB)RemoveOwners(creatorID string,ownerIDs []models.UserID)error{
	stmt, err := db.Prepare(DeleteOwnerByCreator)
	if err!=nil{
		log.Printf("%v",err)
		return database.ErrInternal
	}
	for _,id :=range ownerIDs {
		_, err = stmt.Exec(id.Email, creatorID)
		_,_=db.Query("Update Restaurants set owner_id=null where owner_id=?",id.Email)
		if err != nil {
			log.Printf("%v", err)
			return database.ErrInternal
		}
	}
	return nil
}

func(db *MySqlDB) RemoveOwnersBySuper(ownerIDs []models.UserID)error{
	stmt, err := db.Prepare(DeleteOwner)
	if err!=nil{
		log.Printf("%v",err)
		return database.ErrInternal
	}
	for _,id :=range ownerIDs{
		_,err=stmt.Exec(id.Email)
		_,_=db.Query("Update Restaurants set owner_id=null where owner_id=?",id.Email)
		if err!=nil{
			log.Printf("%v",err)
			return database.ErrInternal
		}
	}
	return nil
}
func(db *MySqlDB)CheckOwnerCreator(creatorID string,ownerIDs []models.UserID)error{
	var creatorIDOut string
	for _,id :=range ownerIDs{
		rows,err:=db.Query("select creator_id from owners where email_id=",id.Email)
		if err!=nil{
			log.Printf("%v",err)
			return database.ErrInternal
		}
		rows.Next()
		err=rows.Scan(&creatorIDOut)
		if err!=nil{
			log.Printf("%v",err)
			return database.ErrInternal
		}
		if creatorIDOut!=creatorID{
			return database.ErrInternal
		}
	}
	return nil

}



//Restaurants functions
func(db *MySqlDB)ShowRestaurantsForSuper()(string,error){
	var result string
	rows,err:=db.Query("select JSON_ARRAYAGG(JSON_OBJECT('id',id,'name', name, 'lat', lat,'lng',lng,'ownerEmailID',owner_id)) from restaurants")
	if err!=nil{
		log.Printf("%v",err)
		return "",database.ErrInternal
	}
	rows.Next()
	rows.Scan(&result)
	return result,nil
}
func(db *MySqlDB)ShowRestaurants(creatorID string)(string,error){
	var result string
	rows,err:=db.Query("select JSON_ARRAYAGG(JSON_OBJECT('id',id,'name', name, 'lat', lat,'lng',lng,'ownerEmailID',owner_id)) from restaurants where creator_id=?",creatorID)
	if err!=nil{
		log.Printf("%v",err)
		return "",database.ErrInternal
	}
	rows.Next()
	rows.Scan(&result)
	return result,nil
}
func(db *MySqlDB)InsertRestaurant(restaurant *models.Restaurant)error{
	stmt,err:=db.Prepare("insert into restaurants(name,lat,lng,creator_id,owner_id) values(?,?,?,?,?)")
	if err!=nil{
		log.Printf("%v",err)
		return database.ErrInternal
	}
	_,err=stmt.Exec(restaurant.Name,restaurant.Lat,restaurant.Lng,restaurant.CreatorID,restaurant.OwnerID)
	if err!=nil{
		log.Printf("%v",err)
		return database.ErrInternal
	}
	return nil
}
func(db *MySqlDB)RemoveRestaurantsBySuper(resIDs []models.ResID)error{
	stmt, err := db.Prepare("delete from restaurants where id=?")
	if err!=nil{
		log.Printf("%v",err)
		return database.ErrInternal
	}
	for _,id :=range resIDs{
		_,err=stmt.Exec(id.ID)
		if err!=nil{
			log.Printf("%v",err)
			return database.ErrInternal
		}
	}
	return nil
}
func(db *MySqlDB)RemoveRestaurants(creatorID string,resIDs []models.ResID)error{
	stmt, err := db.Prepare("delete from restaurants where id=? and creator_id=?")
	if err!=nil{
		log.Printf("%v",err)
		return database.ErrInternal
	}
	for _,id :=range resIDs{
		_,err=stmt.Exec(id.ID,creatorID)
		if err!=nil{
			log.Printf("%v",err)
			return database.ErrInternal
		}
	}
	return nil
}

func(db *MySqlDB)CheckRestaurantCreator(creatorID string,ownerIDs []models.ResID)error{
	var creatorIDOut string
	for _,id :=range ownerIDs{
		rows,err:=db.Query("select creator_id from restaurants where id=?",id.ID)
		if err!=nil{
			log.Printf("%v",err)
			return database.ErrInternal
		}
		rows.Next()
		err=rows.Scan(&creatorIDOut)
		if err!=nil{
			log.Printf("%v",err)
			return database.ErrInternal
		}
		if creatorIDOut!=creatorID{
			return database.ErrInternal
		}
	}
	return nil
}

func(db *MySqlDB)UpdateRestaurant(restaurant *models.RestaurantOutput)error{
	stmt,err:=db.Prepare("update restaurants set name=?,lat=?,lng=?,owner_id=? where id=?")
	if err!=nil{
		log.Printf("%v",err)
		return database.ErrInternal
	}
	isExisting:=db.IsExistingOwner(restaurant.OwnerID)
	if !isExisting{
		return errors.New("owner does not exist try with a different one")
	}
	_,err=stmt.Exec(restaurant.Name,restaurant.Lat,restaurant.Lng,restaurant.OwnerID,restaurant.ID)
	if err!=nil{
			log.Printf("%v",err)
			return database.ErrInternal
	}
	return nil
}



////menu
func(db *MySqlDB)ShowMenu(resID int)(string,error){
	var result string
	rows,err:=db.Query("select JSON_ARRAYAGG(JSON_OBJECT('id',id,'name',name,'price',price)) from dishes where res_id=?",resID)
	if err!=nil{
		log.Printf("%v",err)
		return "",database.ErrInternal
	}
	rows.Next()
	rows.Scan(&result)
	return result,nil
}

func(db *MySqlDB)InsertDishes(dishes []models.Dish,resID int)error{
	stmt,err:=db.Prepare("insert into dishes(res_id,name,price) values(?,?,?)")
	if err!=nil{
		log.Printf("%v",err)
		return database.ErrInternal
	}
	for _,dish:=range dishes{
		_,err=stmt.Exec(resID,dish.Name,dish.Price)
		if err!=nil{
			log.Printf("%v",err)
			return database.ErrInternal
		}
	}
	return nil
}
func(db *MySqlDB)RemoveDishes(ids []models.DishID)error{
	for _,id:=range ids{
		_,err:=db.Query("delete from dishes where id=?",id.ID)
		if err!=nil{
			log.Printf("%v",err)
			return database.ErrInternal
		}
	}
	return nil
}

func(db *MySqlDB)UpdateDishes(dishes []models.DishOutput)error{
	stmt,err:=db.Prepare("update dishes set name=?,price=? where id=?")
	if err!=nil{
		log.Printf("%v",err)
		return database.ErrInternal
	}
	for _,dish:=range dishes{
		_,err=stmt.Exec(dish.Name,dish.Price,dish.ID)
		if err!=nil{
			log.Printf("%v",err)
			return database.ErrInternal
		}
	}
	return nil
}


func(db *MySqlDB)IsExistingOwner(ownerID string) bool{
	count:=0
	rows,err:=db.Query("select count(*) from owners where email_id=?",ownerID)
	rows.Next()
	err=rows.Scan(&count)
	if err!=nil{
		return false
	}
	if count!=1{
		return false
	}
	return true
}


//for owner
func(db *MySqlDB)GetOwnerRestaurants(ownerID string)(string,error){
	var result string
	rows,err:=db.Query("select JSON_ARRAYAGG(JSON_OBJECT('id',id,'name', name, 'lat', lat,'lng',lng)) from restaurants where owner_id=?",ownerID)
	if err!=nil{
		log.Printf("%v",err)
		return "",database.ErrInternal
	}
	rows.Next()
	rows.Scan(&result)
	return result,nil
}
func(db *MySqlDB)CheckRestaurantOwner(ownerID string,resID int)error{
	var ownerIDOut string
	rows,err:=db.Query("select owner_id from restaurants where id=?",resID)
	if err!=nil{
		log.Printf("%v",err)
		return database.ErrInternal
	}
	rows.Next()
	err=rows.Scan(&ownerIDOut)
	if err!=nil{
		log.Printf("%v",err)
		return database.ErrInternal
	}
	if ownerIDOut!=ownerID{
		return database.ErrInternal
	}
	return nil
}

//Helpers for interface functions
func deleteRecordedOwners(db *MySqlDB,owners []models.User){
	stmt,err:=db.Prepare(DeleteOwner)
	for _,owner :=range owners{
		_,err=stmt.Exec(owner.Email)
		if err!=nil{
			log.Printf("error in delRecorded Owners %v",err)
		}
	}
}


