# au
JFrog Artifactory utils

Sample exec: 

au.exe purge --url http://<URL:PORT>/artifactory --repo <ARTIFACTORY_NAME> --user <LOGIN:PASSWORD> -d <NUMBER_DAYS> -r 0
au.exe trash --url http://<URL:PORT>/artifactory --user <LOGIN:PASSWORD> -r 0
au.exe optimize --url http://<URL:PORT>/artifactory --user <LOGIN:PASSWORD> -r 0
