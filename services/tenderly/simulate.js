var apiKey = ''

var obj = JSON.parse(req.body)

var simulateResp = fetch('https://api.tenderly.co/api/v1/account/me/project/project/simulate', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
        'X-Access-Key': apiKey,
    },
    body: JSON.stringify({
        network_id: Number(obj.chain_id || 1),
        block_number: Number(obj.block_number),
        from: obj.from,
        to: obj.to,
        value: obj.value,
        input: obj.input,
        save: true,
        save_if_fails: true,
    }),
})

var simulateRespObj = JSON.parse(simulateResp.body)
if (simulateRespObj.error) {
    throw JSON.stringify(simulateRespObj.error)
}

var simId = simulateRespObj.simulation.id

var shareResp = fetch('https://api.tenderly.co/api/v1/account/me/project/project/simulations/' + simId + '/share', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
        'X-Access-Key': apiKey,
    },
    body: '',
})
_ = shareResp

var outputUrl = 'https://www.tdly.co/shared/simulation/' + simId

JSON.stringify({
    output_url: outputUrl,
    simulation_id: simId,
})
