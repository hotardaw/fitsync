# FitSync

Thanks for taking a look at FitSync. Below are some helpful terminal commands and PGAdmin login credentials/navigation notes.

## Terminal Commands

*Note: all of these commands are called from the "backend" directory (not root!).*

### For All Users

##### Start Docker container
```bash
./../start-app.sh
```

##### Stop Docker container
```bash
docker-compose down -v
```

##### View real-time Docker container logs
```bash
docker logs -f $(docker ps | grep fitsync-backend | awk '{print $1}')
```

### For contributors

##### Compile & generate SQLc code
```bash
docker exec -it $(docker ps | grep fitsync-backend | awk '{print $1}') sqlc generate
```


##### PGAdmin login credentials/navigation guide
If you prefer to see data visualized in table format in the browser, head to <a href="http://localhost:8083/browser/" target="_blank">http://localhost:8083/browser/</a> after starting up the docker container. 

Log in with the first set of credentials:
- `dev@test.com`
- `123lng@#N5las`

Then in the left sidebar, hit the first dropdown (Servers (1)) and log in with the second set of credentials:
- `user01`
- `user01239nTGN35pio!$`

From there, follow absolute path:
Servers (1) > FitSync DB > Databases (2) > fitsyncdb > Schemas (1) > public > Tables (8)