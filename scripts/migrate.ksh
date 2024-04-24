#
# Manage lambda migrations
#

#set -x

# helper definitions
FULL_NAME=$(realpath ${0})
SCRIPT_DIR=$(dirname ${FULL_NAME})
MIGRATE_TOOL=tmp/migrate

function error_and_exit {

   echo "ERROR: ${1}"
   exit 1
}

function checkenv {
   local name=${1}
   local value=${2}

   if [ -z "${value}" ]; then
      error_and_exit "${name} not defined"
   fi
}

# check environment
checkenv "DB_HOST" "${DB_HOST}"
checkenv "DB_PORT" "${DB_PORT}"
checkenv "DB_NAME" "${DB_NAME}"
checkenv "DB_USER" "${DB_USER}"
checkenv "DB_PASSWORD" "${DB_PASSWORD}"

# verify commands
if [ $# -ne 3 ]; then
   echo "use: ${0} <migrate dir> <migrate table> <up|down>"
   exit 1
fi

# for clarity
MIGRATE_DIR=${1}
shift
DB_MIGRATE_TABLE=${1}
shift
MIGRATE_CMD=${1}
shift

# our connection string
CONNECTION_STR="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?x-migrations-table=${DB_MIGRATE_TABLE}"

# run the migration
${MIGRATE_TOOL} -path ${MIGRATE_DIR} --database ${CONNECTION_STR} ${MIGRATE_CMD}

# all over
exit $?

#
# end of file
#
