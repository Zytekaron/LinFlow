#include <iostream>
#include <string>
#include <unordered_set>
#include <chrono>
#include <thread>

#include "kbd.hpp"

extern "C" {
#include <X11/Xlib.h>
#include <X11/Xutil.h>

// TODO CHECK IF NEED (key sym to string fncnlty)
#include <X11/XKBlib.h>
}

const int LINMOD = 269025089;

void start_hook();

std::unordered_set<int> set;
bool linmod_down = false;
char buf[17];

// unused
XComposeStatus comp;
int revert;

void handle_event(
    XEvent ev, Display* d, Window root, Window curFocus, KeySym ks) {
    int len;

    switch (ev.type) {
    case KeyPress:
        len = XLookupString(&ev.xkey, buf, 16, &ks, &comp);
        // prevent held key spam
        if (set.find(ks) != set.end()) break;
        set.insert(ks);

        if (ks == LINMOD) {
            XGrabKeyboard(
                d, curFocus, True, GrabModeAsync, GrabModeAsync, CurrentTime);
            linmod_down = true;

#ifdef DEBUG
            std::cout << "LINMOD pressed\n";
#endif
            break;
        }

        // char sbuf[17];
        // if (len > 0 && isprint(buf[0])) {
        //     buf[len] = 0;
        //     std::sprintf(sbuf, "%s", buf);
        // } else {
        //     std::sprintf(sbuf, "%d", ks);
        // }

        if (linmod_down) {
            native_pass_key(ks, XKeysymToString(ks));
        }

        break;

    case KeyRelease:
        std::this_thread::sleep_for(std::chrono::microseconds(1));
        if (XEventsQueued(d, QueuedAfterReading)) {
            XEvent nev;
            XPeekEvent(d, &nev);

            if (nev.type == KeyPress && nev.xkey.time == ev.xkey.time &&
                nev.xkey.keycode == ev.xkey.keycode) {
                break;
            }
        }

        len = XLookupString(&ev.xkey, buf, 16, &ks, &comp);
        set.erase(ks);

        if (ks == LINMOD) {
            XUngrabKeyboard(d, CurrentTime);
            linmod_down = false;

#ifdef DEBUG
            std::cout << "LINMOD released\n";
#endif

            native_handle_buffer();
        }

        break;

    case FocusOut:
        if (curFocus != root) {
            XSelectInput(d, curFocus, 0);
        }

        XGetInputFocus(d, &curFocus, &revert);
        if (curFocus == PointerRoot) {
            curFocus = root;
        }

        XSelectInput(
            d, curFocus, KeyPressMask | KeyReleaseMask | FocusChangeMask);
    }
}

void start_hook() {
    Display* d = XOpenDisplay(NULL);
    Window root = DefaultRootWindow(d);
    Window curFocus;
    KeySym ks;

    int revert;
    XGetInputFocus(d, &curFocus, &revert);
    XSelectInput(d, curFocus, KeyPressMask | KeyReleaseMask | FocusChangeMask);

    XEvent ev;
    while (true) {
        XNextEvent(d, &ev);

        handle_event(ev, d, root, curFocus, ks);
    }
}
