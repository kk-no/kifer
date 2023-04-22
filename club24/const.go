package club24

const (
	kifuPageURL = "https://www.shogidojo.net/shogi24kifu/"

	loginUserFormXPath     = `//*[@id="uname"]`
	loginPassFormXPath     = `//*[@id="pwd"]`
	loginSubmitButtonXPath = `//*[@id="sub"]`
	
	kifuUser1FormXPath      = `//*[@id="kifusearchTable"]/tbody/tr[1]/td[3]/input[1]`
	kifuUser2FormXPath      = `//*[@id="kifusearchTable"]/tbody/tr[1]/td[3]/input[2]`
	kifuStartDateName       = `fromdate`
	kifuEndDateName         = `todate`
	kifuSearchButtonXPath   = `//*[@id="searchBtn"]`
	kifuDownloadButtonXPath = `//*[@id="dlBtn"]`
)

const (
	dateOnly = "2006/01/02"
)
