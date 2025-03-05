// var runeAPIBase = 'http://host.docker.internal:8580' // beta
var runeAPIBase = 'http://host.docker.internal:8570' // prod

var chainNameMap = {
    1: 'ethereum',
    56: 'bsc',
}

function firstOr(x, defaultValue) {
    if (x && x.length > 0) {
        return x[0]
    }
    return defaultValue
}

function getAmountByFunc(chainId, addr, start, end) {
    var funcName = 'FUNGetUserMintedBrBtcAmountALLCHAIN'
    var params = {
        user: addr,
        start_time: start,
        end_time: end,
    }

    if (chainId != 0 && chainNameMap[chainId] != '') {
        funcName = 'FUNGetUserMintedBrBtcAmountCHAIN'
        params.chain_name = chainNameMap[chainId]
    }

    var payload = {
        func_name: funcName,
        params: JSON.stringify(params)
    }

    var resp = fetch(runeAPIBase + '/dsn/execsql', {
        method: 'POST',
        body: JSON.stringify(payload)
    })

    var respObj = JSON.parse(resp.body)
    var mintedAmount = 0
    if (respObj.data && respObj.data[0]) {
        mintedAmount = Number(respObj.data[0].total_amount)
    }
    return mintedAmount
}

var chainId = firstOr(req.query.chain_id, 0)
var addr = firstOr(req.query.address, '')
var start = firstOr(req.query.from, 0)
var end = firstOr(req.query.to, 0)
var amountLimit = firstOr(req.query.amount, 0)

var brBTCAmount = getAmountByFunc(chainId, addr, start, end)

JSON.stringify({
    result: brBTCAmount >= amountLimit
})
