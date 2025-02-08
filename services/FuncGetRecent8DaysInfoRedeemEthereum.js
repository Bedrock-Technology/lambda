// var rune_api_base = 'http://host.docker.internal:8580' // beta
var rune_api_base = 'http://host.docker.internal:8570' // prod

var payload = {
    func_name: 'FuncGetRecent8DaysInfoRedeemEthereum',
    params: JSON.stringify({
        token: req.query.tokens[0].split(',')
    })
}

var resp = fetch(rune_api_base + '/rune/execsql', {
    method: 'POST',
    body: JSON.stringify(payload)
})

resp.body
