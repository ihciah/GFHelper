all: goblin ork register packetdecryptor orkmre empadder dbmanager dbmigrator usergenerator

goblin:
	go build -o bin/Goblin Goblin.go

ork:
	go build -o bin/Ork Ork.go

register:
	go build -o bin/Register Register.go

packetdecryptor:
	go build -o bin/PacketDecryptor PacketDecryptor.go

orkmre:
	go build -o bin/OrkMRE OrkMRE.go

empadder:
	go build -o bin/EmpAdder EmpAdder.go

dbmanager:
	go build -o bin/DBManager DBManager.go

dbmigrator:
	go build -o bin/DBMigrator DBMigrator.go

usergenerator:
	go build -o bin/UserGenerator UserGenerator.go

.PHONY: all goblin ork register packetdecryptor orkmre empadder dbmanager dbmigrator usergenerator