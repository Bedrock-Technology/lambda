var { vars, net } = LambdaHelper

var apiBase = vars.merkle_api_base
var req = vars.req

function firstOr(x, defaultValue) {
    if (x && x.length > 0) {
        return x[0]
    }
    return defaultValue
}

var resp = net.fetch(apiBase+ '/api/v1/merkle/proof', {
    method: 'POST',
    body: JSON.stringify({
        address: firstOr(req.query.address, ''),
    })
})

resp.body

