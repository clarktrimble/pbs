
TKO="travelpic/JuneJuly2023"
TKO="takeout"
RZD="travelpic-resized/resized"
RZD="takeout-resized"

ROOT="/Users/trimble"

export PBL_TAKEOUTPATH="$ROOT/$TKO"
export PBL_RESIZEDPATH="$ROOT/$RZD"
export PBL_FILTER='PXL_20230[67]'
export PBL_APICLIENT_BASEURI="http://localhost:8088"
export PBL_TRUNCATE=99
export PBL_DRYRUN="true"

export PBAPI_BOLT_PATH="photo.db"
export PBAPI_TRUNCATE=99
export PBAPI_SERVER_PORT=8088
export PBAPI_SERVER_TIMEOUT="30s"

export PBR_TAKEOUTPATH="$ROOT/$TKO"
export PBR_RESIZEDPATH="$ROOT/$RZD"
export PBR_FILTER='PXL_20230[67]'
export PBR_TRUNCATE=99
export PBR_DRYRUN="true"

## Todo: think about using the same prefix for all yeah?
