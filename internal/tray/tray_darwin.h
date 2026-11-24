#include <stdbool.h>
#include <stdint.h>

void tray_run(void);
void tray_quit(void);
void tray_set_title(const char *title);
void tray_add_item(uint32_t id, const char *title, const char *tooltip, bool checked);
void tray_add_separator(void);
void tray_set_item_title(uint32_t id, const char *title);
void tray_set_item_checked(uint32_t id, bool checked);
void tray_set_item_hidden(uint32_t id, bool hidden);
