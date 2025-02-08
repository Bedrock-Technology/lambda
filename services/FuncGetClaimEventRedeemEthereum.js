// var rune_api_base = 'http://host.docker.internal:8580' // beta
var rune_api_base = 'http://host.docker.internal:8570' // prod

var payload = {
    func_name: 'FuncGetClaimEventRedeemEthereum',
    params: JSON.stringify({
        offset: req.query.offset[0],
        limit: req.query.limit[0],
        start_time: req.query.start_time[0],
        end_time: req.query.end_time[0],
        user: req.query.user[0],
    })
}

var resp = fetch(rune_api_base + '/rune/execsql', {
    method: 'POST',
    body: JSON.stringify(payload)
})

var obj = JSON.parse(resp.body)
obj.data = JSON.parse(obj.data[0].result)
JSON.stringify(obj)
