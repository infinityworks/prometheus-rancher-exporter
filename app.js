var request = require('request')
var debug = {
    log: require('debug')('rs'),
    http: require('debug')('rs.http')
}
var async = require('async')
var promclient = require('prometheus-client')

process.on('SIGTERM', function () {
    process.exit(0);
});

var opts = getOptions()
createServer(opts.host, opts.port, opts.listen_port, opts.update_interval)

function getOptions() {
    var opts = {
        host:            process.env.HOST,
        port:            process.env.PORT || 8080,
        api_access_key:  process.env.API_ACCESS_KEY,
        api_secret_key:  process.env.API_SECRET_KEY,
        listen_port:     process.env.LISTEN_PORT || 9010,
        update_interval: process.env.UPDATE_INTERVAL || 5000
    }

    var requiredOpts = [
        'HOST',
        'API_ACCESS_KEY',
        'API_SECRET_KEY'
    ]
    requiredOpts.forEach(function(name) {
        if (!opts[name.toLowerCase()]) {
            throw new Error('Missing environment variable for option: ' + name)
            process.exit(1)
        }
    })

    return opts
}

function createServer(host, port, listen_port, update_interval) {
    var client = new promclient()
    var gauges = {}

    function createGauge(name) {
        gauges[name] = client.newGauge({
            namespace: 'rancher_status',
            name: name,
            help: '1 if all containers in an environment are active'
        })
    }

    function updateGauge(name, value) {
        if (!gauges[name]) {
            createGauge(name)
        }
        gauges[name].set({
            name: name
        }, value)
    }

    function updateMetrics() {
        debug.log('requesting metrics')
        getEnvironmentsState(host, port, function(err, results) {
            if (err) {
                throw err
            }
            debug.log('got metric results %o', results)
            Object.keys(results).forEach(function(name) {
                var state = results[name]
                name = name.replace('-', '_')
                updateGauge(name, state == 'active' ? 1 : 0)
            });
        });
    }

    client.listen(listen_port)
    debug.log('listening on %s', listen_port)

    updateMetrics()
    setInterval(updateMetrics, update_interval)
}


function getEnvironmentsState(host, port, callback) {
    var envIdMap = {}

    async.waterfall([
        function(next) {
            var uri = 'http://' + host + ':' + port + '/v1/projects'
            jsonRequest(uri, function(err, json) {
                if (err) {
                    return next(err)
                }
                var environments = json.data[0].links.environments
                next(null, environments)
            })
        },
        function(environmentsUrl, next) {
            jsonRequest(environmentsUrl, function(err, json) {
                if (err) {
                    return next(err)
                }
                var servicesUrl = json.data.map(function(raw) {
                    return raw.links.services
                });
                json.data.forEach(function(env) {
                    envIdMap[env.id] = env.name
                });
                next(null, servicesUrl)
            });
        },
        function(servicesUrls, next) {
            var tasks = servicesUrls.map(function(servicesUrl) {
                return function(next) {
                    jsonRequest(servicesUrl, next)
                }
            });

            async.parallel(tasks, function(err, results) {
                var data = results.map(function(servicesRaw) {
                    return servicesRaw.data
                });

                next(null, data)
            });
        },
        function(servicesData, next) {
            var services = servicesData.map(function(stackServices) {
                return stackServices.map(function(service) {
                    return {
                        name: service.name,
                        state: service.state,
                        environment: envIdMap[service.environmentId]
                    }
                });
            });

            var flattened = []
            services.forEach(function(service) {
                flattened = flattened.concat(service)
            });

            next(null, flattened)
        },
        function(serviceData, next) {
            var envState = {}
            serviceData.forEach(function(service) {
                if (!envState[service.environment]) {
                    envState[service.environment] = service.state
                } else if (service.state != 'active') {
                    envState[service.environment] = service.state
                }
            });
            next(null, envState)
        }
    ], function(err, results) {
        callback(err, results)
    })
}

function jsonRequest(uri, callback) {
    debug.http('send request: %s', uri)

    request({
        uri: uri,
        headers: {
            'Accept': 'application/json'
        },
        auth: {
            user: opts.api_access_key,
            pass: opts.api_secret_key,
            sendImmediately: true
        }
    }, function(err, response, body) {
        if (err) {
            debug.http('got error response')
            return callback(err)
        }

        debug.http('got valid response with status %s', response.statusCode)

        var data;
        try {
            data = JSON.parse(body)
        } catch(e) {
            debug.http('could not JSON decode response body')
            var error = new Error('could not decode')
            error.response = response
            error.body = body
            return callback(error)
        }

        return callback(null, data)
    })
}
