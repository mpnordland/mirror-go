Floating Window Support
=======================
Floating windows will require some modification to the Leaf struct.
Mainly a new field of type *xrect.XRect. This new field will indicate
that this leaf has a user defined position/size. A new method will be
added to the Node specification called Floating() if this method returns
true, then that Node's Layout method will be called with the XRect passed to its parent's
Layout function and space will not be allocated for it. Then that Node's Layout method will be responsible
for positioning itself if, it has a user defined pos/size, it will place itself to that.
otherwise, it will center itself according to the XRect passed to it.

Modal window management
=======================
Mirror uses a modal style of
window management which is very like
vim. There will be three modes avalilible:
    -Command Mode
        You enter this mode by pressing Mod4
        this mode captures all key presses and mouse clicks
        this mode allows you to communicate with the window manager
        and do actions like focus windows, launch apps, and close windows
    -Interact Mode
        you enter this mode by pressing 'i' in Command mode
        this mode allows keypresses and mouse clicks
        to pass through to applications
    -Visual Mode
        you enter this mode by pressing 'v' in Command mode
        Like command mode, this mode also captures all key presses and mouse clicks
        this mode allows you to select multiple windows by 
        focusing them and pressing enter
        There are a few commands also such as 's' to swap two windows
