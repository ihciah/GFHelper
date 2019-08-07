echo "Building..."
go build -o bin/Goblin Goblin.go
go build -o bin/Ork Ork.go
go build -o bin/Register Register.go
go build -o bin/PacketDecryptor PacketDecryptor.go
go build -o bin/OrkMRE OrkMRE.go
go build -o bin/EmpAdder EmpAdder.go
go build -o bin/DBManager DBManager.go
go build -o bin/DBMigrator DBMigrator.go
go build -o bin/UserGenerator UserGenerator.go
echo "Done."