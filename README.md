
# mongodb
brew services start mongodb/brew/mongodb-community
brew services stop mongodb/brew/mongodb-community
brew services restart mongodb/brew/mongodb-community
brew services list 

# mongoDB terminal command for checking data
mongosh
show dbs
use [name of database]
show collections
db.getCollection('match_logs').find()
