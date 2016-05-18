var mongojs = require('mongojs');
var devicesDb = require('../database').collection('devices');
var uuid = require('node-uuid');
var config = require('../../config');
var log = require('../logger');
var os = require('os');


/** createDevice() */
exports.createDevice = function(req, cb) {
    var res = {};
    var deviceId = uuid.v1();

    /** TODO: req must contain deviceId, deviceName and deviceType - form validation */
    var entity = {
        'deviceId' : deviceId,
        'deviceName' : req.deviceName,
        'deviceType' : req.deviceType
    }
        
    /** Save the device and check for errors */
    devicesDb.save(entity, function(err, device) {
        if (err) {
            res = JSON.stringify({
                'status' : 505,
                'message' : 'Internal Server Error'
            });
            return cb(res, err);
        }

        res = JSON.stringify({
                'status': 200,
                'message': 'OK',
                'deviceId': deviceId
        });

        return cb(err, res);
    });
}

/** getDevices() */
exports.getDevices = function(req, cb) {
    var res = {};

    devicesDb.find(function(err, devices) {
        if (err) {
            res = JSON.stringify({
                'status' : 505,
                'message': 'Internal Server Error'
            });
            return cb(err, res);
        }

        console.log("DEVICES: ", devices);

        res = JSON.stringify(devices);
        return cb(err, res);
    });
}

/** getDevice() */
exports.getDevice = function(req, cb) {
    var res = {};

    devicesDb.findOne({'deviceId': req.deviceId}, function(err, device) {
        if (err) {
            res = JSON.stringify({
                'status' : 505,
                'message': 'Internal Server Error'
            });
            return cb(err, res);
        }
        
        if (device) {
            res = JSON.stringify(device);
        } else {
            res = res = JSON.stringify({
                'status' : 404,
                'message': 'Not found',
                'deviceId': req.deviceId
            });
        }
        return cb(err, res);
    });
}

/** updateDevice() */
exports.updateDevice = function(req, cb) {
    var res = {};

    /** Use our device model to find the device we want */
    console.log(req);

    devicesDb.update(
            {deviceId: req.deviceId},
            {$set: req},
            function(err, upd) {
            if (err) {
                res = JSON.stringify({
                    'status' : 505,
                    'message': 'Internal Server Error'
                });
                return cb(err, res);
            }

            if (upd.n === 1) {
                res = JSON.stringify({
                    'status': 200,
                    'message': 'OK',
                    'deviceId': req.deviceId
                });
            } else {
                res = JSON.stringify({
                    'status' : 404,
                    'message': 'Not found',
                    'deviceId': req.deviceId
                });
            }

            return cb(err, res);
    });
}

/** deleteDevice() */
exports.deleteDevice = function(req, cb) {
    var res = {};

    devicesDb.remove({
        deviceId: req.deviceId
    }, function(err, rem) {
        if (err) {
            res = JSON.stringify({
                'status' : 505,
                'message': 'Internal Server Error'
            });
            return cb(err, res);
        }

        if (rem.n === 1) {
            res = JSON.stringify({
                'status': 200,
                'message': 'OK',
                'deviceId': req.deviceId
            });
        } else {
            res = JSON.stringify({
                'status': 404,
                'message': 'Not Found',
                'deviceId': req.deviceId
            });
        }

        return cb(err, res);
    });
}

