package cmd

//
// import (
// 	//"fmt"
//
// 	"github.com/codegangsta/cli"
// 	//"github.com/containerops/vessel/models"
// )
//
// var CmdDatabase = cli.Command{
// 	Name:        "db",
// 	Usage:       "Database init, backup and maintain",
// 	Description: "Develop environment use tidb and goleveldb, production please use tidb and HBase.",
// 	Action:      runDatabase,
// 	Flags: []cli.Flag{
// 		cli.StringFlag{
// 			Name:  "action",
// 			Usage: "Action，[sync/backup/restore]",
// 		},
// 		cli.BoolFlag{
// 			Name:  "verbose",
// 			Usage: "export verbose when exec sync action",
// 		},
// 		cli.BoolFlag{
// 			Name:  "force",
// 			Usage: "delete all tables when exec sync action",
// 		},
// 	},
// }
//
// func runDatabase(c *cli.Context) {
// 	/*	if len(c.String("action")) > 0 {
// 			action := c.String("action")
//
// 			switch action {
// 			case "sync":
// 				if err := models.Sync(c.Bool("force"), c.Bool("verbose")); err != nil {
// 					fmt.Println("Init database struct error, ", err.Error())
// 				}
// 				break
// 			default:
// 				break
// 			}
// 		}
// 	*/
// }
