#import <AppKit/AppKit.h>
#import "spotify_darwin.h"

static NSAppleEventDescriptor *spotifyTarget(void) {
    NSRunningApplication *app = [NSRunningApplication
        runningApplicationsWithBundleIdentifier:@"com.spotify.client"].firstObject;
    // Address the existing PID so a poll cannot launch Spotify, even if it quits.
    return app ? [NSAppleEventDescriptor descriptorWithProcessIdentifier:app.processIdentifier] : nil;
}

static NSAppleEventDescriptor *sendEvent(NSAppleEventDescriptor *target, AEEventClass eventClass,
                                         AEEventID eventID, NSAppleEventDescriptor *object,
                                         NSString **message) {
    NSAppleEventDescriptor *event = [NSAppleEventDescriptor
        appleEventWithEventClass:eventClass eventID:eventID targetDescriptor:target
        returnID:kAutoGenerateReturnID transactionID:kAnyTransactionID];
    if (object) {
        [event setParamDescriptor:object forKeyword:keyDirectObject];
    }

    NSError *error = nil;
    NSAppleEventDescriptor *reply = [event sendEventWithOptions:NSAppleEventSendDefaultOptions
        timeout:2 error:&error];
    if (!reply) {
        *message = error.localizedDescription ?: @"Spotify did not return an Apple Event reply";
        return nil;
    }
    SInt32 code = [reply paramDescriptorForKeyword:keyErrorNumber].int32Value;
    if (code != noErr) {
        if (code == errAEEventNotPermitted) {
            *message = @"Spotify automation permission denied. Allow access in System Settings > Privacy & Security > Automation.";
            return nil;
        }
        NSString *detail = [reply paramDescriptorForKeyword:keyErrorString].stringValue;
        *message = [NSString stringWithFormat:@"Spotify Apple Event error %d: %@",
                    (int)code, detail ?: @"no description"];
        return nil;
    }
    return reply;
}

static NSAppleEventDescriptor *propertyObject(OSType property, NSAppleEventDescriptor *container) {
    NSAppleEventDescriptor *object = [NSAppleEventDescriptor recordDescriptor];
    [object setDescriptor:[NSAppleEventDescriptor descriptorWithTypeCode:typeProperty]
               forKeyword:keyAEDesiredClass];
    [object setDescriptor:[NSAppleEventDescriptor descriptorWithEnumCode:formPropertyID]
               forKeyword:keyAEKeyForm];
    [object setDescriptor:[NSAppleEventDescriptor descriptorWithTypeCode:property]
               forKeyword:keyAEKeyData];
    [object setDescriptor:container ?: [NSAppleEventDescriptor nullDescriptor]
               forKeyword:keyAEContainer];
    return [object coerceToDescriptorType:typeObjectSpecifier];
}

static NSAppleEventDescriptor *readProperty(NSAppleEventDescriptor *target, OSType property,
                                            NSAppleEventDescriptor *container, NSString **message) {
    if (*message) {
        return nil;
    }
    NSAppleEventDescriptor *reply = sendEvent(target, kAECoreSuite, kAEGetData,
                                             propertyObject(property, container), message);
    if (!reply) {
        return nil;
    }
    NSAppleEventDescriptor *value = [reply paramDescriptorForKeyword:keyDirectObject];
    if (!value) {
        *message = @"Spotify returned no property value";
    }
    return value;
}

static NSString *readText(NSAppleEventDescriptor *target, OSType property,
                           NSAppleEventDescriptor *container, NSString **message) {
    NSAppleEventDescriptor *value = readProperty(target, property, container, message);
    NSString *text = value.stringValue;
    if (value && !text) {
        *message = @"Spotify returned a non-text track property";
    }
    return text;
}

static double readNumber(NSAppleEventDescriptor *target, OSType property,
                          NSAppleEventDescriptor *container, NSString **message) {
    NSAppleEventDescriptor *value = readProperty(target, property, container, message);
    NSAppleEventDescriptor *number = [value coerceToDescriptorType:typeIEEE64BitFloatingPoint];
    if (value && !number) {
        *message = @"Spotify returned a non-numeric playback property";
    }
    return number.doubleValue;
}

char *spotify_read(bool progress, bool album, char **error) {
    @autoreleasepool {
        NSAppleEventDescriptor *target = spotifyTarget();
        if (!target) {
            return strdup("{}");
        }

        NSString *message = nil;
        OSType state = readProperty(target, 'pPlS', nil, &message).enumCodeValue;
        NSMutableDictionary *data = [NSMutableDictionary dictionary];
        switch (state) {
            case 'kPSP': data[@"PlayerState"] = @"playing"; break;
            case 'kPSp': data[@"PlayerState"] = @"paused"; break;
            case 'kPSS': data[@"PlayerState"] = @"stopped"; break;
            default:
                if (!message) {
                    message = @"Spotify returned an unknown player state";
                }
        }

        if (!message && state != 'kPSS') {
            NSAppleEventDescriptor *track = propertyObject('pTrk', nil);
            NSString *name = readText(target, 'pnam', track, &message);
            NSString *artist = readText(target, 'pArt', track, &message);
            if (name) data[@"Track"] = name;
            if (artist) data[@"Artist"] = artist;
            if (album) {
                NSString *title = readText(target, 'pAlb', track, &message);
                if (title) data[@"Album"] = title;
            }
            if (progress) {
                data[@"Duration"] = @(readNumber(target, 'pDur', track, &message) / 1000.0);
                data[@"Position"] = @(readNumber(target, 'pPos', nil, &message));
            }
        }
        if (message) {
            *error = strdup(message.UTF8String);
            return NULL;
        }

        NSError *jsonError = nil;
        NSData *json = [NSJSONSerialization dataWithJSONObject:data options:0 error:&jsonError];
        if (!json) {
            *error = strdup(jsonError.localizedDescription.UTF8String);
            return NULL;
        }
        NSString *text = [[NSString alloc] initWithData:json encoding:NSUTF8StringEncoding];
        return strdup(text.UTF8String);
    }
}

char *spotify_toggle(void) {
    @autoreleasepool {
        NSAppleEventDescriptor *target = spotifyTarget();
        if (!target) {
            return strdup("Spotify is not running");
        }
        NSString *message = nil;
        sendEvent(target, 'spfy', 'PlPs', nil, &message);
        return message ? strdup(message.UTF8String) : NULL;
    }
}
