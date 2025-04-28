var { vars, utils, crypto } = LambdaHelper

function firstOr(x, defaultValue) {
    if (x && x.length > 0) {
        return x[0]
    }
    return defaultValue
}

function require(x, msg) {
    if (!x) {
        throw JSON.stringify({ msg: msg })
    }
}

var camps = {
    'bera_superweek_0414': ['0xa645b923FC178881f85b5658Bae6dAE0C24D7390', 'assets/bera_0424.csv', 1],
}

var campId = firstOr(vars.req.query.campaign)
require(camps[campId], 'Invalid campaign')

var addr = firstOr(vars.req.query.addr, '')
var standardAddr = utils.hex_to_address(addr)

var amounts = {}
var records = utils.csv_read(camps[campId][1])
for (var r of records) {
    require(r.length == 2, 'Invalid CSV format')
    amounts[utils.hex_to_address(r[0])] = r[1]
}

require(amounts[standardAddr], 'Address is not found')

var tree = crypto.merkle(records.map(r => {
    return {
        address: r[0],
        amount: r[1],
    }
}))

var resp = {
    contract: camps[campId][0],
    root: tree.root,
    epoch: camps[campId][2],
    addr: standardAddr,
    amount: amounts[standardAddr],
    proof: tree.proof[standardAddr],
}

JSON.stringify({
    code: 200,
    msg: 'ok',
    data: resp,
})
