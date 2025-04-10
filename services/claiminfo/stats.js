var { db } = LambdaHelper

function require(x, msg) {
    if (!x) {
        throw JSON.stringify({ msg: msg })
    }
}

var data = db.select(`select * from public.pre_deposit_stats`)

JSON.stringify({
    code: 200,
    data: data,
})
