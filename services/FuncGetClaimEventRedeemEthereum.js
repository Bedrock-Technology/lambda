var { vars, net } = LambdaHelper

var runeAPIBase = vars.rune_api_base
var req = vars.req

function firstOr(x, defaultValue) {
    if (x && x.length > 0) {
        return x[0]
    }
    return defaultValue
}

var payload = {
    func_name: 'FuncGetClaimEventRedeemEthereum',
    params: JSON.stringify({
        offset: firstOr(req.query.offset, ''),
        limit: firstOr(req.query.limit, ''),
        start_time: firstOr(req.query.start_time, ''),
        end_time: firstOr(req.query.end_time, ''),
        user: firstOr(req.query.user, ''),
    })
}

var resp = net.fetch(runeAPIBase + '/dsn/execsql', {
    method: 'POST',
    body: JSON.stringify(payload)
})

var obj = JSON.parse(resp.body)
obj.data = JSON.parse(obj.data[0].result)
JSON.stringify(obj)
