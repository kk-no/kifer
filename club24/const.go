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

	// File format is following.
	// 開始日時：yyyy/mm/dd hh:mm:ss
	// 終了日時：yyyy/mm/dd hh:mm:ss
	// 棋戦：game rule
	// 手合割：handicap
	// 先手：user1(9999)
	// 後手：user2(9999)
	// Moves...
	// 64 投了
	// まで63手で先手の勝ち
	fileMetadataCount = 6 + 2
)

const (
	dateOnly = "2006/01/02"
)
