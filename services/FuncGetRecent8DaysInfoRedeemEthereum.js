var { vars, net } = LambdaHelper

var runeAPIBase = vars.rune_api_base
var req = vars.req

var token = []
if (req.query.tokens && req.query.tokens.length > 0) {
    token = req.query.tokens[0].split(',')
}

var payload = {
    func_name: 'FuncGetRecent8DaysInfoRedeemEthereum',
    params: JSON.stringify({ token: token })
}

var resp = net.fetch(rune_api_base + '/dsn/execsql', {
    method: 'POST',
    body: JSON.stringify(payload)
})

resp.body
