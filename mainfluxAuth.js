/**
 * Copyright (c) Mainflux
 *
 * Mainflux server is licensed under an Apache license, version 2.0 license.
 * All rights not explicitly granted in the Apache license, version 2.0 are reserved.
 * See the included LICENSE file for more details.
 */

var config = require('./config');
var deviceController = require('./app/controllers/devices');
var nats = require('nats').connect('nats://' + config.nats.host + ':' + config.nats.port);

var replySubject = '';

/**
 * createDevice()
 */
var createDevice = function(req) {
    console.log("createDevice");
    
    deviceController.createDevice(req, function(err, res) {

        console.log("RESULT: ", res)
        /** Send reply to the API server */
        nats.publish(replySubject, res);
    });
}

/**
 * getDevices()
 */
var getDevices = function(req) {
    console.log("getDevices");
 
    deviceController.getDevices(req, function(err, res) {

        console.log("RESULT: ", res)
        console.log(typeof res);

        nats.publish(replySubject, res);
    });
}

/**
 * getDevice()
 */
var getDevice = function(req) {
    console.log("getDevice");

    deviceController.getDevice(req, function(err, res) {
        nats.publish(replySubject, res);
    });
}

/**
 * updateDevice()
 */
var updateDevice = function(req) {
    console.log("updateDevice");

    deviceController.updateDevice(req, function(err, res) {

        /** Reply to the caller */
        nats.publish(replySubject, res);

        /** Fan out the update to other services */
        nats.publish('core_out', res);

    });
}

/**
 * deleteDevice()
 */
var deleteDevice = function(req) {
    console.log("deleteDevice");

    deviceController.deleteDevice(req, function(err, res) {
        nats.publish(replySubject, res);
    });
}


/**
 * Lookup table of all API functions
 */
var fnList = {
    'createDevice' : createDevice,
    'getDevices' : getDevices,
    'getDevice' : getDevice,
    'updateDevice' : updateDevice,
    'deleteDevice' : deleteDevice
}

/** Subscribe to core_in subject */
nats.subscribe('core_in', function(req, replyTo) {
    console.log('Received a message: ' + req);

    var rpc = JSON.parse(req);
    replySubject = replyTo;

    console.log(replySubject);
    
    /** Call the function */
    fnList[rpc.method](rpc.body)

});


var banner = `
                                     
_|      _|            _|                _|_|  _|                      
_|_|  _|_|    _|_|_|      _|_|_|      _|      _|  _|    _|  _|    _|  
_|  _|  _|  _|    _|  _|  _|    _|  _|_|_|_|  _|  _|    _|    _|_|    
_|      _|  _|    _|  _|  _|    _|    _|      _|  _|    _|  _|    _|  
_|      _|    _|_|_|  _|  _|    _|    _|      _|    _|_|_|  _|    _|  
                                                                      
    
                == Industrial IoT System ==
       

                Made with <3 by Mainflux Team

[w] http://mainflux.io
[t] @mainflux

`

console.log(banner);

