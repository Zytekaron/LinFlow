#include <iostream>
#include <string>
#include <unordered_set>
#include <chrono>
#include <thread>

#include "kbd.hpp"

extern "C" {
#include <X11/Xlib.h>
#include <X11/Xutil.h>
}

int LINMOD = -1;

std::unordered_set<int> set;
bool linmod_down = false;
char buf[17];

// unused
XComposeStatus comp;
int revert;

void handle_event(XEvent event, Display* display, Window root, Window currentFocus, KeySym keySym) {
    int len;

    switch (event.type) {
    case KeyPress:
        len = XLookupString(&event.xkey, buf, 16, &keySym, &comp);
        // prevent held key spam
        if (set.find(keySym) != set.end()) {
            break;
        }
        set.insert(keySym);

        if (keySym == LINMOD) {
            XGrabKeyboard(display, currentFocus, True, GrabModeAsync, GrabModeAsync, CurrentTime);
            linmod_down = true;
            break;
        }

        if (linmod_down) {
            native_pass_key(keySym, XKeysymToString(keySym));
        }

        break;

    case KeyRelease:
        std::this_thread::sleep_for(std::chrono::microseconds(1));
        if (XEventsQueued(display, QueuedAfterReading)) {
            XEvent nextEvent;
            XPeekEvent(display, &nextEvent);

            if (nextEvent.type == KeyPress && nextEvent.xkey.time == event.xkey.time
                && nextEvent.xkey.keycode == event.xkey.keycode) {
                break;
            }
        }

        len = XLookupString(&event.xkey, buf, 16, &keySym, &comp);
        set.erase(keySym);

        if (keySym == LINMOD) {
            XUngrabKeyboard(display, CurrentTime);
            linmod_down = false;
            native_handle_buffer();
        }

        break;

    case FocusOut:
        if (currentFocus != root) {
            XSelectInput(display, currentFocus, 0);
        }

        XGetInputFocus(display, &currentFocus, &revert);
        if (currentFocus == PointerRoot) {
            currentFocus = root;
        }

        XSelectInput(display, currentFocus, KeyPressMask | KeyReleaseMask | FocusChangeMask);
    }
}

void start_hook(int key) {
    LINMOD = key;
    Display* d = XOpenDisplay(NULL);
    Window currentFocus, root = DefaultRootWindow(d);
    KeySym ks;

    int revert;
    XGetInputFocus(d, &currentFocus, &revert);
    XSelectInput(d, currentFocus, KeyPressMask | KeyReleaseMask | FocusChangeMask);

    XEvent ev;
    while (true) {
        XNextEvent(d, &ev);
        handle_event(ev, d, root, currentFocus, ks);
    }
}
