package shared

var mapQuery = map[ServiceName]string{
	ASIPCNT: `
			SELECT 
			  timestamp, 
			  message, 
			  flowid, 
			  type, 
			  hostname, 
			  part, 
			  context.correlationid, 
			  context.event, 
			  context.uri, 
			  context.providerid, 
			  context.providerhotelid, 
			  context.locale, 
			  context.providerbrandid, 
			  context.providerchainid 
			FROM 
			  "s3log"."pcntapirqrs_init"
			`,
	ASIPRSV: ``,
	ASIPSRC: `
			SELECT 
			  timestamp, 
			  message, 
			  flowid, 
			  type, 
			  hostname, 
			  part, 
			  context.correlationid, 
			  context.providerid, 
			  context.event, 
			  context.sourcemarket, 
			  context.checkindate, 
			  context.checkoutdate, 
			  context.locale, 
			  context.currency,
			  context.noofadult,
			  context.noofchild,
			  context.noofroom
			FROM 
			  "s3log"."psrclogrqrs_init"
			`,
	ASIPSRC_DETAIL: `
			SELECT 
			  timestamp, 
			  message, 
			  flowid, 
			  type, 
			  hostname, 
			  part, 
			  context.correlationid
			FROM 
			  "s3log"."psrclogdetail_init"
			`,
}
