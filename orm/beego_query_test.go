package orm

/*
run these SQL before test

CREATE DATABASE IF NOT EXISTS `beego_query`;

USE `beego_query`;

DROP TABLE IF EXISTS `Company`;

CREATE TABLE `Company` (
  `id` varchar(10) NOT NULL,
  `name` varchar(50) NOT NULL,
  `capital` decimal(12,0) NOT NULL,
  `established` date NOT NULL,
  `headquarter` varchar(20) NOT NULL,
  PRIMARY KEY (`id`)
)

INSERT INTO `Company` VALUES
  ('BM','Buns Motor Ltd.',3800000000,'1952-01-01','Munich'),
  ('OC','Orange Computer Inc.',6000000000,'1983-01-01','San Francisco'),
  ('PE','Pony Electronics Ltd.',1500000000,'1975-01-01','Tokyo'),
  ('RB','Roys Bank Corp.',4200000000,'1931-01-01','London'),
  ('YRP','Yellow River Properties Ltd.',2300000000,'1960-01-01','Hong Kong');
*/

import (
	"os"
	"testing"
	"time"

	beego "github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type Company struct {
	Id          string    `orm:"pk;column(id)"`
	Name        string    `orm:"column(name)"`
	Capital     int64     `orm:"column(capital)"`
	Established time.Time `orm:"column(established);type(date)"`
	Headquarter string    `orm:"column(headquarter)"`
}

var ormer beego.Ormer

func TestMain(m *testing.M) {
	beego.RegisterModel(new(Company))
	beego.RegisterDataBase("default", "mysql", "root:pass@tcp(localhost:3306)/beego_query?charset=utf8&parseTime=true")
	//beego.Debug = true
	ormer = beego.NewOrm()
	os.Exit(m.Run())
}

func TestSingleOp1(t *testing.T) {
	var company Company
	err := ormer.QueryTable("Company").
		Filter("Name", "Pony Electronics Ltd.").
		One(&company)
	if err != nil {
		t.Fatal(err)
	}
	if company.Id != "PE" {
		t.Fatalf("id expected=%s; actual=%s", "PE", company.Id)
	}
}

func TestSingleOp2(t *testing.T) {
	var companies []*Company
	hkt, _ := time.LoadLocation("Asia/Shanghai")
	date := time.Date(1970, time.January, 1, 0, 0, 0, 0, hkt)

	num, err := ormer.QueryTable("Company").
		Filter("Established__gte", date).
		OrderBy("Established").
		All(&companies)
	if err != nil {
		t.Fatal(err)
	}
	if num != 2 {
		t.Fatalf("num expected=%d; actual=%d", 2, num)
	}
	if companies[0].Id != "PE" {
		t.Fatalf("1st id expected=%s; actual=%s", "PE", companies[0].Id)
	}
	if companies[1].Id != "OC" {
		t.Fatalf("2nd id expected=%s; actual=%s", "OC", companies[1].Id)
	}
}

func TestSingleOp3(t *testing.T) {
	var companies []*Company

	num, err := ormer.QueryTable("Company").
		Filter("Capital__lte", 3800000000).
		OrderBy("Capital").
		All(&companies)
	if err != nil {
		t.Fatal(err)
	}
	if num != 3 {
		t.Fatalf("num expected=%d; actual=%d", 3, num)
	}
	if companies[0].Id != "PE" {
		t.Fatalf("1st id expected=%s; actual=%s", "PE", companies[0].Id)
	}
	if companies[1].Id != "YRP" {
		t.Fatalf("2nd id expected=%s; actual=%s", "YRP", companies[1].Id)
	}
	if companies[2].Id != "BM" {
		t.Fatalf("3rd id expected=%s; actual=%s", "BM", companies[2].Id)
	}
}

func TestSingleOp4(t *testing.T) {
	var companies []*Company

	num, err := ormer.QueryTable("Company").
		Filter("Headquarter__in", "San Francisco", "London").
		OrderBy("Id").
		All(&companies)
	if err != nil {
		t.Fatal(err)
	}
	if num != 2 {
		t.Fatalf("num expected=%d; actual=%d", 2, num)
	}
	if companies[0].Id != "OC" {
		t.Fatalf("1st id expected=%s; actual=%s", "OC", companies[0].Id)
	}
	if companies[1].Id != "RB" {
		t.Fatalf("2nd id expected=%s; actual=%s", "RB", companies[1].Id)
	}
}

func TestAnd(t *testing.T) {
	var companies []*Company
	num, err := ormer.QueryTable("Company").
		Filter("Name__contains", "u").
		Filter("Headquarter__contains", "n").
		OrderBy("Id").
		All(&companies)
	if err != nil {
		t.Fatal(err)
	}
	if num != 2 {
		t.Fatalf("num expected=%d; actual=%d", 2, num)
	}
	if companies[0].Id != "BM" {
		t.Fatalf("1st id expected=%s; actual=%s", "BM", companies[0].Id)
	}
	if companies[1].Id != "OC" {
		t.Fatalf("2nd id expected=%s; actual=%s", "OC", companies[1].Id)
	}
}

func TestOr(t *testing.T) {
	var companies []*Company
	num, err := ormer.QueryTable("Company").
		SetCond(beego.Or(
			beego.Cond("Name__startswith", "Ro"),
			beego.Cond("Headquarter__endswith", "ng"))).
		OrderBy("Id").
		All(&companies)
	if err != nil {
		t.Fatal(err)
	}
	if num != 2 {
		t.Fatalf("num expected=%d; actual=%d", 2, num)
	}
	if companies[0].Id != "RB" {
		t.Fatalf("1st id expected=%s; actual=%s", "RB", companies[0].Id)
	}
	if companies[1].Id != "YRP" {
		t.Fatalf("2nd id expected=%s; actual=%s", "YRP", companies[1].Id)
	}
}

func TestNot(t *testing.T) {
	var companies []*Company
	num, err := ormer.QueryTable("Company").
		Exclude("Name__contains", "Ltd.").
		OrderBy("Id").
		All(&companies)
	if err != nil {
		t.Fatal(err)
	}
	if num != 2 {
		t.Fatalf("num expected=%d; actual=%d", 2, num)
	}
	if companies[0].Id != "OC" {
		t.Fatalf("1st id expected=%s; actual=%s", "OC", companies[0].Id)
	}
	if companies[1].Id != "RB" {
		t.Fatalf("2nd id expected=%s; actual=%s", "RB", companies[1].Id)
	}
}

func TestAndNot(t *testing.T) {
	var companies []*Company
	num, err := ormer.QueryTable("Company").
		Filter("Headquarter__contains", "n").
		Exclude("Name__contains", "B").
		OrderBy("Id").
		All(&companies)
	if err != nil {
		t.Fatal(err)
	}
	if num != 2 {
		t.Fatalf("num expected=%d; actual=%d", 2, num)
	}
	if companies[0].Id != "OC" {
		t.Fatalf("1st id expected=%s; actual=%s", "OC", companies[0].Id)
	}
	if companies[1].Id != "YRP" {
		t.Fatalf("2nd id expected=%s; actual=%s", "YRP", companies[1].Id)
	}
}

func TestAndOr(t *testing.T) {
	var companies []*Company
	hkt, _ := time.LoadLocation("Asia/Shanghai")
	date := time.Date(1980, time.January, 1, 0, 0, 0, 0, hkt)

	and := beego.And(
			beego.Cond("Capital__gt", 4000000000),
			beego.Cond("Established__lt", date))
	cond := beego.Or(
			beego.Cond("Headquarter", "Munich"),
			and)
	num, err := ormer.QueryTable("Company").
		SetCond(cond).
		OrderBy("Id").
		All(&companies)
	if err != nil {
		t.Fatal(err)
	}
	if num != 2 {
		t.Fatalf("num expected=%d; actual=%d", 2, num)
	}
	if companies[0].Id != "BM" {
		t.Fatalf("1st id expected=%s; actual=%s", "BM", companies[0].Id)
	}
	if companies[1].Id != "RB" {
		t.Fatalf("2nd id expected=%s; actual=%s", "RB", companies[1].Id)
	}
}

func TestOrNotAnd(t *testing.T) {
	var companies []*Company
	hkt, _ := time.LoadLocation("Asia/Shanghai")
	date := time.Date(1950, time.January, 1, 0, 0, 0, 0, hkt)
	
	not := beego.Not(
			beego.Cond("Name__contains", "Ltd."))
	or := beego.Or(
			beego.Cond("Headquarter__startswith", "To"),
			not)
	cond := beego.And(
			beego.Cond("Established__gt", date),
			or)
	num, err := ormer.QueryTable("Company").
		SetCond(cond).
		OrderBy("Id").
		All(&companies)
	if err != nil {
		t.Fatal(err)
	}
	if num != 2 {
		t.Fatalf("num expected=%d; actual=%d", 2, num)
	}
	if companies[0].Id != "OC" {
		t.Fatalf("1st id expected=%s; actual=%s", "OC", companies[0].Id)
	}
	if companies[1].Id != "PE" {
		t.Fatalf("2nd id expected=%s; actual=%s", "PE", companies[1].Id)
	}
}
