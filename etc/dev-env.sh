
TKO="travelpic/JuneJuly2023"
TKO="takeout"

RZD="travelpic-resized/resized"
RZD="takeout-resized"

ROOT="/Users/trimble"

export PB_TAKEOUTPATH="$ROOT/$TKO"
export PB_RESIZEDPATH="$ROOT/$RZD"
export PB_MATCH='PXL_20230[67]'
export PB_TRUNCATE=99
export PB_DRYRUN="true"

export PB_APICLIENT_BASEURI="http://localhost:8088"

export PB_BOLT_PATH="photo.db"
export PB_SERVER_PORT=8088
export PB_SERVER_TIMEOUT="30s"

