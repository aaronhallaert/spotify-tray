#import <AppKit/AppKit.h>
#import "tray_darwin.h"

extern void tray_ready(void);
extern void tray_item_clicked(uint32_t id);

@interface SpotifyTrayDelegate : NSObject <NSApplicationDelegate>
@property(strong) NSStatusItem *statusItem;
@property(strong) NSMenu *menu;
@end

@implementation SpotifyTrayDelegate

- (void)applicationDidFinishLaunching:(NSNotification *)notification {
    self.statusItem = [[NSStatusBar systemStatusBar] statusItemWithLength:NSVariableStatusItemLength];
    self.menu = [[NSMenu alloc] init];
    self.menu.autoenablesItems = NO;
    self.statusItem.menu = self.menu;
    tray_ready();
}

- (void)itemClicked:(NSMenuItem *)item {
    tray_item_clicked((uint32_t)item.tag);
}

@end

static SpotifyTrayDelegate *delegate;

static void onMain(void (^work)(void)) {
    if ([NSThread isMainThread]) {
        work();
    } else {
        dispatch_sync(dispatch_get_main_queue(), work);
    }
}

void tray_run(void) {
    @autoreleasepool {
        NSApplication *app = [NSApplication sharedApplication];
        delegate = [[SpotifyTrayDelegate alloc] init];
        app.delegate = delegate;
        [app run];
    }
}

void tray_quit(void) {
    dispatch_async(dispatch_get_main_queue(), ^{
        [NSApp stop:nil];
        NSEvent *event = [NSEvent otherEventWithType:NSEventTypeApplicationDefined
                                           location:NSZeroPoint
                                      modifierFlags:0
                                          timestamp:0
                                       windowNumber:0
                                            context:nil
                                            subtype:0
                                              data1:0
                                              data2:0];
        [NSApp postEvent:event atStart:NO];
    });
}

void tray_set_title(const char *title) {
    @autoreleasepool {
        NSString *text = [NSString stringWithUTF8String:title];
        onMain(^{
            delegate.statusItem.button.title = text;
        });
    }
}

void tray_add_item(uint32_t itemID, const char *title, const char *tooltip, bool checked) {
    @autoreleasepool {
        NSString *text = [NSString stringWithUTF8String:title];
        NSString *hint = [NSString stringWithUTF8String:tooltip];
        onMain(^{
            NSMenuItem *item = [delegate.menu addItemWithTitle:text
                                                       action:@selector(itemClicked:)
                                                keyEquivalent:@""];
            item.target = delegate;
            item.tag = itemID;
            item.toolTip = hint;
            item.state = checked ? NSControlStateValueOn : NSControlStateValueOff;
        });
    }
}

void tray_add_separator(void) {
    onMain(^{
        [delegate.menu addItem:[NSMenuItem separatorItem]];
    });
}

void tray_set_item_title(uint32_t itemID, const char *title) {
    @autoreleasepool {
        NSString *text = [NSString stringWithUTF8String:title];
        onMain(^{
            [delegate.menu itemWithTag:itemID].title = text;
        });
    }
}

void tray_set_item_checked(uint32_t itemID, bool checked) {
    onMain(^{
        [delegate.menu itemWithTag:itemID].state = checked ? NSControlStateValueOn : NSControlStateValueOff;
    });
}

void tray_set_item_hidden(uint32_t itemID, bool hidden) {
    onMain(^{
        [delegate.menu itemWithTag:itemID].hidden = hidden;
    });
}
