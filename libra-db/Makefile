MIGRATION_DIR ='../migrations'
MIGRATE_CMD=cd bin && migrate -path $(MIGRATION_DIR) \
	-database 'postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)'

build: migrate-up

migrate-%: setup
	$(MIGRATE_CMD) $*

new: setup
	cd bin && migrate create -ext sql -dir $(MIGRATION_DIR) $(NAME)

setup: check_env
	mkdir -p bin
	cd bin && test -f migrate || curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
	chmod +x ./bin/migrate

clean:
	rm -rf bin

check_env:
ifeq ($(DB_USER),)
$(error missing required environment DB_USER)
else ifeq ($(DB_PASSWORD),)
$(error missing required environment DB_PASSWORD)
else ifeq ($(DB_HOST),)
$(error missing required environment DB_HOST)
else ifeq ($(DB_PORT),)
$(error missing required environment DB_PORT)
else ifeq ($(DB_NAME),)
$(error missing required environment DB_NAME)
endif
