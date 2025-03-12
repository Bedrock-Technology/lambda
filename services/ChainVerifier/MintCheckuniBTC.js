var { vars, net } = LambdaHelper

var runeAPIBase = vars.rune_api_base
var req = vars.req

var chainNameMap = {
    1: 'ethereum',
    10: 'optimism',
    56: 'bsc',
    5000: 'mantle',
    34443: 'mode',
    42161: 'arbitrum',
    223: 'b2',
    80094: 'bera',
    200901: 'bitlayer',
    60808: 'bob',
    4200: 'merlin',
    7000: 'zeta'
}

function firstOr(x, defaultValue) {
    if (x && x.length > 0) {
        return x[0]
    }
    return defaultValue
}

function getAmountByFunc(funcNames, chainId, addr, start, end) {
    var funcName = funcNames[0]
    var params = {
        user: addr,
        start_time: start,
        end_time: end,
    }

    if (chainId != 0 && chainNameMap[chainId] != '') {
        funcName = funcNames[1]
        params.chain_name = chainNameMap[chainId]
    }

    var payload = {
        func_name: funcName,
        params: JSON.stringify(params)
    }

    var resp = net.fetch(runeAPIBase + '/dsn/execsql', {
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

var uniBTCAmount = getAmountByFunc(['FUNGetUserMintedUniBtcAmountALLCHAIN', 'FUNGetUserMintedUniBtcAmountCHAIN'], chainId, addr, start, end)

JSON.stringify({
    result: uniBTCAmount >= amountLimit
})
