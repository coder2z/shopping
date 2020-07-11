package test

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"testing"
)

var (
	MysqlHandler *gorm.DB
)

func mysqlBuild() *gorm.DB {
	var err error
	DB, err := gorm.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			"root",
			"123456",
			"127.0.0.1:3306",
			"test"))

	if err != nil {
		log.Panicf("models.Setup err: %v", err)
		return nil
	}

	//	设置表前缀
	//gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
	//	return setting.DatabaseSetting.TablePrefix + defaultTableName
	//}

	DB.SingularTable(true)
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)
	return DB
}

// `Department` 属于 `User`， 外键是`UserID`
type Department struct {
	gorm.Model
	UserID uint
	User   User
	Name   string
}

type User struct {
	gorm.Model
	Name string
}

func TestGormBelongsTo(t *testing.T) {
	MysqlHandler = mysqlBuild()
	MysqlHandler.AutoMigrate(User{})
	MysqlHandler.AutoMigrate(Department{})

	//多表链表添加数据
	//profile := Profile{
	//	User: User{
	//		Name: "杨泽淼",
	//	},
	//	Name: "项目组",
	//}
	//MysqlHandler.Create(&profile)

	//多表链表查询
	//var infoList []Department
	//var user User
	//MysqlHandler.Preload("User").Find(&infoList)
	//for _, value := range infoList {
	//	fmt.Println(value.Name)
	//	fmt.Println(value.User.Name)
	//}

	// 查询添加条件 这里的条件就是 Department id为2 对应的数据
	//var info Department
	//info.ID = 2
	//MysqlHandler.Debug().Preload("User").Find(&info)
	//fmt.Println(info.Name, info.User.Name)

	//USER name为杨泽淼A 对应的数据
	//var info Department
	//MysqlHandler.Debug().Preload("User", func(query *gorm.DB) *gorm.DB {
	//	return query.Where("name =? ", "杨泽淼A")
	//}).First(&info)
	//fmt.Println(info.Name, info.User.Name, info.User.ID)
	//	or
	//var infoList []Department
	//MysqlHandler.Debug().Preload("User", "name =?","杨泽淼A").First(&infoList)
	//for _, value := range infoList {
	//	fmt.Println(value.Name, value.User.Name, value.User.ID)
	//}

	// 使用 Related 查找 belongs to 关系
	var department Department               //需要查找总的结构体
	department.ID = 3                       //定义查询条件
	MysqlHandler.Debug().First(&department) //首先查询总表
	/// SELECT * FROM `department`  WHERE `department`.`deleted_at` IS NULL AND `department`.`id` = 3 ORDER BY `department`.`id` ASC LIMIT 1
	MysqlHandler.Debug().Model(&department).Related(&department.User) //查询子表进行赋值
	/// SELECT * FROM `user`  WHERE `user`.`deleted_at` IS NULL AND ((`id` = 2))
	fmt.Println(department.Name, "---", department.User.Name, "---", department.ID, "---", department.User.ID)

}

//车
type Car struct {
	gorm.Model
	Host         string //车主人名字
	LicensePlate LicensePlate
}

//车牌
type LicensePlate struct {
	gorm.Model
	Number string //车牌号
	CarID  uint
}

func TestGormHasOne(t *testing.T) {
	MysqlHandler = mysqlBuild()
	MysqlHandler.AutoMigrate(Car{})
	MysqlHandler.AutoMigrate(LicensePlate{})

	//多表链表添加数据
	//car := Car{
	//	Host: "李四B",
	//	LicensePlate: LicensePlate{
	//		Number: "川A666666",
	//	},
	//}
	//	INSERT INTO `car` (`created_at`,`updated_at`,`deleted_at`,`host`) VALUES ('2020-07-11 10:30:53','2020-07-11 10:30:53',NULL,'李四A')/
	//	INSERT INTO `license_plate` (`created_at`,`updated_at`,`deleted_at`,`number`,`car_id`) VALUES ('2020-07-11 10:30:53','2020-07-11 10:30:53',NULL,'川A123456',1)
	//MysqlHandler.Debug().Create(&car)

	//查询 Preload
	//var car []Car
	//MysqlHandler.Debug().Preload("LicensePlate").Find(&car)
	//	SELECT * FROM `car`  WHERE `car`.`deleted_at` IS NULL
	//	SELECT * FROM `license_plate`  WHERE `license_plate`.`deleted_at` IS NULL AND ((`car_id` IN (1,2)))
	//for _, value := range car {
	//	fmt.Println(value.Host, value.LicensePlate.Number)
	//}

	//条件查询 Preload 这种方法只适合查询单条数据，而且在主表就确认了只有一条 ,不然可能会造成查询出来的结构体数据很多垃圾数据，需要自己手动清洗
	var car Car
	MysqlHandler.Debug().Where("Host=?", "李四A").Preload("LicensePlate", func(query *gorm.DB) *gorm.DB {
		return query.Where("number =? ", "川A123456")
	}).Find(&car)
	fmt.Println(car.Host, car.LicensePlate.Number)

	// 使用 Related 查找 belongs to 关系
	//var car Car               //需要查找总的结构体
	//car.ID = 2                       //定义查询条件
	//MysqlHandler.Debug().First(&car) //首先查询总表
	//MysqlHandler.Debug().Model(&car).Related(&car.LicensePlate) //查询子表进行赋值
	//fmt.Println(car.Host, car.LicensePlate.Number)
}

//学校
type School struct {
	gorm.Model
	Name       string
	Profession []Profession
}

//专业
type Profession struct {
	gorm.Model
	Name     string
	SchoolId uint
}

func TestGormHasMany(t *testing.T) {
	MysqlHandler = mysqlBuild()
	MysqlHandler.AutoMigrate(School{})
	MysqlHandler.AutoMigrate(Profession{})

	//创建
	//profession1 := Profession{Name: "信息工程"}
	//profession2 := Profession{Name: "计算机科学"}
	//var professionList []Profession
	//professionList = append(professionList, profession1)
	//professionList = append(professionList, profession2)
	//school := School{
	//	Name:       "成都大学",
	//	Profession: professionList,
	//}
	//MysqlHandler.Save(&school) //or Create

	//查询单列
	//var school School
	//MysqlHandler.Preload("Profession").First(&school)
	//fmt.Println(school)

	//var school School
	//MysqlHandler.First(&school)
	//MysqlHandler.Model(&school).Related(&school.Profession)
	//fmt.Println(school)

	//查询多列
	var school []School
	//MysqlHandler.First(&school)
	MysqlHandler.Model(&school).Preload("Profession").Find(&school)
	for _, value := range school {
		fmt.Println(value.Name, value.Profession)
	}
}

type Users struct {
	gorm.Model
	Name      string
	Languages []Language `gorm:"many2many:user_languages;"`
}

type Language struct {
	gorm.Model
	Name  string
	Users []Users `gorm:"many2many:user_languages;"`
}

func TestGormManyToMany(t *testing.T) {
	MysqlHandler = mysqlBuild()
	MysqlHandler.AutoMigrate(Users{})
	MysqlHandler.AutoMigrate(Language{})

	//	创建
	//langEN := Language{Name: "EN"}
	//langCN := Language{Name: "CN"}
	//u1 := &Users{
	//	Name: "user1",
	//	Languages: []Language{
	//		langEN,
	//		langCN,
	//	},
	//}
	//MysqlHandler.Create(u1)

	//LangCN := Language{}
	//MysqlHandler.Where("name=?","CN").First(&LangCN)
	//u2 := &Users{
	//	Name: "user2",
	//	Languages: []Language{
	//		LangCN,
	//	},
	//}
	//MysqlHandler.Create(u2)

	//查询
	//获取 用户id 为 3 的 user 的语言：
	//var users Users
	//MysqlHandler.Find(&users, 3)
	//MysqlHandler.Model(&users).Related(&users.Languages, "Languages")


	//查询
	//获取 使用语言 为 CN 的 user：
	//var language Language
	//MysqlHandler.Find(&language, "name=?","CN")
	//MysqlHandler.Model(&language).Related(&language.Users, "Users")
	//fmt.Println(language)

}
