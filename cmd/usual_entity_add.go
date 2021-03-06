package cmd

import (
	"github.com/fatih/color"
	"gopkg.in/abiosoft/ishell.v2"
	"strings"
)

const replaceCommentLink = `//generator insert entity`
const replaceCommentObject = `//generator insert object_entity`

func usualEntityAdd(c *ishell.Context) {

	yellow := color.New(color.FgYellow).SprintFunc()
	c.Println(yellow("Start creating new entity and api"))

	WoDbModel, _ := GetOsArgument(WithoutDbModels.ToString())
	Uuid, _ := GetOsArgument(UuidAsPk.ToString())

	entity, err := getEntity(c)

	if err != nil {
		return
	}

	CamelCase := strings.Title(entity)
	snakeCase := getLowerCase(entity)
	firstLowerCase := GetFirstLowerCase(entity)

	sourceFile := "./settings/routes.go"
	destinationFile := "./settings/routes.go"

	CopyFile(
		sourceFile,
		destinationFile,
		[]string{"\n// route-constant-generator here dont touch this line", "{Entity}", "{entity}"},
		[]string{usualTemplateSettingsRoutesConstEntity.Content, CamelCase, firstLowerCase},
		c)

	sourceFile = "./settings/routes.go"
	destinationFile = "./settings/routes.go"

	CopyFile(
		sourceFile,
		destinationFile,
		[]string{"\n    // router-list-generator here dont touch this line", "{Entity}", "{entity}"},
		[]string{usualTemplateSettingsRoutesListEntity.Content, CamelCase, firstLowerCase},
		c)

	sourceFile = "./router/router.go"
	destinationFile = "./router/router.go"

	CopyFile(
		sourceFile,
		destinationFile,
		[]string{"router-generator here dont touch this line", "{Entity}", "{entity}"},
		[]string{getRouteContent(), CamelCase, firstLowerCase},
		c)


	if ! WoDbModel.BoolResult {

		sourceFile = "./logic/assigner.go"
		destinationFile = "./logic/assigner.go"

		CopyFile(
			sourceFile,
			destinationFile,
			[]string{"// add all assign functions", "{Entity}", "{entity}"},
			[]string{getAssignContent(), CamelCase, firstLowerCase},
			c)
	}

	CreateFileIfNotExists(usualTemplateGen.Path, usualTemplateGen.Content, c)

	sourceFile = "./generator/" + snakeCase + ".go"
	destinationFile = "./generator/" + snakeCase + ".go"
	CreateFile(sourceFile, getEntityGenContent(), c)
	CopyFile(
		sourceFile,
		destinationFile,
		[]string{"{entity-name}", "{Entity}", "{entity}"},
		[]string{CamelCase, CamelCase, firstLowerCase},
		c)

	sourceFile = "./webapp/" + snakeCase + ".go"
	destinationFile = "./webapp/" + snakeCase + ".go"
	CreateFile(sourceFile, getWebAppContent(), c)

	CopyFile(
		sourceFile,
		destinationFile,
		[]string{"{entity-name}", "{Entity}", "{entity}"},
		[]string{CamelCase, CamelCase, firstLowerCase},
		c)

	switch CamelCase {
	case "Auth", "User":
		break
	default:
		sourceFile = "./webapp/" + snakeCase + "_test.go"
		destinationFile = "./webapp/" + snakeCase + "_test.go"

		CreateFile(sourceFile, getWebAppTestContent(), c)

		CopyFile(
			sourceFile,
			destinationFile,
			[]string{"{entity-name}", "{Entity}", "{entity}"},
			[]string{CamelCase, CamelCase, firstLowerCase},
			c)
		break
	}

	sourceFile = "./types/" + snakeCase + ".go"
	destinationFile = "./types/" + snakeCase + ".go"

	CreateFile(sourceFile, getTypeContent(), c)

	CopyFile(
		sourceFile,
		destinationFile,
		[]string{"{entity-name}", "{Entity}", "{entity}"},
		[]string{CamelCase, CamelCase, firstLowerCase},
		c)


	if ! WoDbModel.BoolResult {

		sourceFile = "./dbmodels/" + snakeCase + ".go"
		destinationFile = "./dbmodels/" + snakeCase + ".go"

		CreateFile(sourceFile, getDbModelContent(Uuid.BoolResult), c)

		CopyFile(
			sourceFile,
			destinationFile,
			[]string{"{entity-name}", "{Entity}", "{entity}"},
			[]string{CamelCase, CamelCase, firstLowerCase},
			c)
	}

	sourceFile = "./logic/" + snakeCase + ".go"
	destinationFile = "./logic/" + snakeCase + ".go"

	CreateFile(sourceFile, getLogicContent(), c)

	CopyFile(
		sourceFile,
		destinationFile,
		[]string{"{entity-name}", "{Entity}", "{entity}"},
		[]string{CamelCase, CamelCase, firstLowerCase},
		c)

	if ! WoDbModel.BoolResult {

		sourceFile = "./bootstrap/insert_data_to_db.go"
		destinationFile = "./bootstrap/insert_data_to_db.go"

		replaceFrom := replaceCommentLink
		replaceTo := replaceCommentLink + `
          ` + "&dbmodels." + CamelCase + "{},"

		CopyFile(
			sourceFile,
			destinationFile,
			[]string{replaceFrom},
			[]string{replaceTo},
			c)

		replaceFrom = replaceCommentObject
		replaceTo = replaceCommentObject + `
          ` + "dbmodels." + CamelCase + "{},"

		CopyFile(
			sourceFile,
			destinationFile,
			[]string{replaceFrom},
			[]string{replaceTo},
			c)
	}

	c.Println("New entity " + CamelCase + " was created")
}

func getLogicContent() (c string) {

	// DISABLED crudParams. Just comment in ROUTES for routes not in logic

	//crudArgs, _ := GetOsArgument("crud")
	crudParams := Crud{}

	WoDbModel, _ := GetOsArgument(WithoutDbModels.ToString())

	//if len(crudArgs.StringResult) > 0 {
	//
	//	crudParams.IsFind = strings.Contains(crudArgs.StringResult, "f")
	//	crudParams.IsCreate = strings.Contains(crudArgs.StringResult, "c")
	//	crudParams.IsRead = strings.Contains(crudArgs.StringResult, "r")
	//	crudParams.IsUpdate = strings.Contains(crudArgs.StringResult, "u")
	//	crudParams.IsDelete = strings.Contains(crudArgs.StringResult, "d")
	//	crudParams.IsFindOrCreate = strings.Contains(crudArgs.StringResult, "a")
	//	crudParams.IsUpdateOrCreate = strings.Contains(crudArgs.StringResult, "x")
	//} else {
		crudParams = Crud{true, true, true, true, true, true, true}
	//}

	c = GetUsualTemplateLogicContent(crudParams, WoDbModel.BoolResult)
	return
}

func getAssignContent() (c string) {

	cfg := TypeConfig{}

	hasId, _ := GetOsArgument("no-id")
	cfg.IsId = !hasId.BoolResult

	c = GetUsualTemplateAssignContent(cfg)
	return
}

func getTypeContent() (c string) {

	cfg := TypeConfig{}

	hasId, _ := GetOsArgument("no-id")
	cfg.IsId = !hasId.BoolResult

	IsUuid, _ := GetOsArgument(UuidAsPk.ToString())
	cfg.IsUuid = IsUuid.BoolResult

	c = GetUsualTemplateTypeContent(cfg)
	return
}

func getWebAppContent() (webappContent string) {

	webappContent = usualTemplateWebappEntity.Content

	AuthcrudArgs, _ := GetOsArgument(CheckAuth.ToString())
	authParams := Crud{true, true, true, true, true, true, true}

	if len(AuthcrudArgs.StringResult) > 0 {

		authParams.IsFind = strings.Contains(AuthcrudArgs.StringResult, "f")
		authParams.IsCreate = strings.Contains(AuthcrudArgs.StringResult, "c")
		authParams.IsRead = strings.Contains(AuthcrudArgs.StringResult, "r")
		authParams.IsUpdate = strings.Contains(AuthcrudArgs.StringResult, "u")
		authParams.IsDelete = strings.Contains(AuthcrudArgs.StringResult, "d")
		authParams.IsFindOrCreate = strings.Contains(AuthcrudArgs.StringResult, "a")
		authParams.IsUpdateOrCreate = strings.Contains(AuthcrudArgs.StringResult, "x")

	}

	// DISABLED crudParams. Just comment in ROUTES for routes not in logic

	//methodCrudArgs, _ := GetOsArgument("crud")
	methodCrudParams := Crud{true, true, true, true, true, true, true}
	//
	//if len(methodCrudArgs.StringResult) > 0 {
	//
	//	methodCrudParams.IsFind = strings.Contains(methodCrudArgs.StringResult, "f")
	//	methodCrudParams.IsCreate = strings.Contains(methodCrudArgs.StringResult, "c")
	//	methodCrudParams.IsRead = strings.Contains(methodCrudArgs.StringResult, "r")
	//	methodCrudParams.IsUpdate = strings.Contains(methodCrudArgs.StringResult, "u")
	//	methodCrudParams.IsDelete = strings.Contains(methodCrudArgs.StringResult, "d")
	//	methodCrudParams.IsFindOrCreate = strings.Contains(methodCrudArgs.StringResult, "a")
	//	methodCrudParams.IsUpdateOrCreate = strings.Contains(methodCrudArgs.StringResult, "x")
	//
	//}

	webappContent = assignMsName(GetUsualTemplateWebAppContent(authParams, methodCrudParams))

	return
}

func getEntityGenContent() (genContent string) {

	genContent = usualTemplateGenEntity.Content

	return
}

func getWebAppTestContent() (webappTestContent string) {

	webappTestContent = usualTemplateWebappTestEntity.Content

	//methodCrudArgs, _ := GetOsArgument("crud")
	methodCrudParams := Crud{true, true, true, true, true, true, true}
	//
	//if len(methodCrudArgs.StringResult) > 0 {
	//
	//	methodCrudParams.IsFind = strings.Contains(methodCrudArgs.StringResult, "f")
	//	methodCrudParams.IsCreate = strings.Contains(methodCrudArgs.StringResult, "c")
	//	methodCrudParams.IsRead = strings.Contains(methodCrudArgs.StringResult, "r")
	//	methodCrudParams.IsUpdate = strings.Contains(methodCrudArgs.StringResult, "u")
	//	methodCrudParams.IsDelete = strings.Contains(methodCrudArgs.StringResult, "d")
	//	methodCrudParams.IsFindOrCreate = strings.Contains(methodCrudArgs.StringResult, "a")
	//	methodCrudParams.IsUpdateOrCreate = strings.Contains(methodCrudArgs.StringResult, "x")
	//
	//}

	webappTestContent = assignMsName(GetUsualTemplateWebAppTestContent(methodCrudParams))

	return
}

func getEntity(c *ishell.Context) (entity string, err error) {

	var arguments RegularFind

	arguments, err = GetOsArgument("entity")

	if len(arguments.StringResult) < 1 || err != nil {
		return getName(c, false, "Entity")
	}

	entity = arguments.StringResult

	return
}

func getRouteContent() string {

	crudParams := Crud{
		IsFind:         true,
		IsCreate:       true,
		IsRead:         true,
		IsUpdate:       true,
		IsDelete:       true,
		IsFindOrCreate: true,
		IsUpdateOrCreate: true,
	}

	crudArgs, _ := GetOsArgument("crud")

	if len(crudArgs.StringResult) > 0 {

		crudParams.IsFind = strings.Contains(crudArgs.StringResult, "f")
		crudParams.IsCreate = strings.Contains(crudArgs.StringResult, "c")
		crudParams.IsRead = strings.Contains(crudArgs.StringResult, "r")
		crudParams.IsUpdate = strings.Contains(crudArgs.StringResult, "u")
		crudParams.IsDelete = strings.Contains(crudArgs.StringResult, "d")
		crudParams.IsFindOrCreate = strings.Contains(crudArgs.StringResult, "a")
		crudParams.IsUpdateOrCreate = strings.Contains(crudArgs.StringResult, "x")

	}

	return GetUsualTemplateRouteEntity(crudParams)
}
