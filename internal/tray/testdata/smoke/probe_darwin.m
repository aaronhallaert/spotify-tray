#import <AppKit/AppKit.h>

int smoke_check(int checked) {
    __block int result = 0;
    dispatch_sync(dispatch_get_main_queue(), ^{
        NSStatusItem *status = [(NSObject *)NSApp.delegate valueForKey:@"statusItem"];
        NSMenu *menu = status.menu;
        if (![status.button.title isEqualToString:@"♫ demo"]) result = 1;
        else if (menu.numberOfItems != 3) result = 2;
        else if (![[menu itemAtIndex:0].title isEqualToString:@"Setting"]) result = 3;
        else if ([menu itemAtIndex:0].state != checked) result = 4;
        else if (![menu itemAtIndex:1].separatorItem) result = 5;
        else if (![[menu itemAtIndex:2].title isEqualToString:@"Pause"]) result = 6;
        else if ([menu itemAtIndex:2].hidden) result = 7;
    });
    return result;
}

int smoke_hidden(void) {
    __block BOOL hidden;
    dispatch_sync(dispatch_get_main_queue(), ^{
        NSStatusItem *status = [(NSObject *)NSApp.delegate valueForKey:@"statusItem"];
        hidden = [status.menu itemAtIndex:2].hidden;
    });
    return hidden;
}

void smoke_click(void) {
    dispatch_sync(dispatch_get_main_queue(), ^{
        NSStatusItem *status = [(NSObject *)NSApp.delegate valueForKey:@"statusItem"];
        NSMenuItem *item = [status.menu itemAtIndex:0];
        [NSApp sendAction:item.action to:item.target from:item];
    });
}
