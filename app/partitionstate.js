// source: app/partition_state.proto
/**
 * @fileoverview
 * @enhanceable
 * @suppress {messageConventions} JS Compiler reports an error if a variable or
 *     field starts with 'MSG_' and isn't a translatable message.
 * @public
 */
// GENERATED CODE -- DO NOT EDIT!

goog.provide('proto.PartitionState');

goog.require('jspb.BinaryReader');
goog.require('jspb.BinaryWriter');
goog.require('jspb.Message');

/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.PartitionState = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.PartitionState.repeatedFields_, null);
};
goog.inherits(proto.PartitionState, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.PartitionState.displayName = 'proto.PartitionState';
}

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.PartitionState.repeatedFields_ = [3];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.PartitionState.prototype.toObject = function(opt_includeInstance) {
  return proto.PartitionState.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.PartitionState} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.PartitionState.toObject = function(includeInstance, msg) {
  var f, obj = {
    cumulativeTimesteps: jspb.Message.getFloatingPointFieldWithDefault(msg, 1, 0.0),
    partitionIndex: jspb.Message.getFieldWithDefault(msg, 2, 0),
    stateList: (f = jspb.Message.getRepeatedFloatingPointField(msg, 3)) == null ? undefined : f
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.PartitionState}
 */
proto.PartitionState.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.PartitionState;
  return proto.PartitionState.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.PartitionState} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.PartitionState}
 */
proto.PartitionState.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readDouble());
      msg.setCumulativeTimesteps(value);
      break;
    case 2:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setPartitionIndex(value);
      break;
    case 3:
      var value = /** @type {!Array<number>} */ (reader.readPackedDouble());
      msg.setStateList(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.PartitionState.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.PartitionState.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.PartitionState} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.PartitionState.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getCumulativeTimesteps();
  if (f !== 0.0) {
    writer.writeDouble(
      1,
      f
    );
  }
  f = message.getPartitionIndex();
  if (f !== 0) {
    writer.writeInt64(
      2,
      f
    );
  }
  f = message.getStateList();
  if (f.length > 0) {
    writer.writePackedDouble(
      3,
      f
    );
  }
};


/**
 * optional double cumulative_timesteps = 1;
 * @return {number}
 */
proto.PartitionState.prototype.getCumulativeTimesteps = function() {
  return /** @type {number} */ (jspb.Message.getFloatingPointFieldWithDefault(this, 1, 0.0));
};


/**
 * @param {number} value
 * @return {!proto.PartitionState} returns this
 */
proto.PartitionState.prototype.setCumulativeTimesteps = function(value) {
  return jspb.Message.setProto3FloatField(this, 1, value);
};


/**
 * optional int64 partition_index = 2;
 * @return {number}
 */
proto.PartitionState.prototype.getPartitionIndex = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {number} value
 * @return {!proto.PartitionState} returns this
 */
proto.PartitionState.prototype.setPartitionIndex = function(value) {
  return jspb.Message.setProto3IntField(this, 2, value);
};


/**
 * repeated double state = 3;
 * @return {!Array<number>}
 */
proto.PartitionState.prototype.getStateList = function() {
  return /** @type {!Array<number>} */ (jspb.Message.getRepeatedFloatingPointField(this, 3));
};


/**
 * @param {!Array<number>} value
 * @return {!proto.PartitionState} returns this
 */
proto.PartitionState.prototype.setStateList = function(value) {
  return jspb.Message.setField(this, 3, value || []);
};


/**
 * @param {number} value
 * @param {number=} opt_index
 * @return {!proto.PartitionState} returns this
 */
proto.PartitionState.prototype.addState = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 3, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.PartitionState} returns this
 */
proto.PartitionState.prototype.clearStateList = function() {
  return this.setStateList([]);
};


