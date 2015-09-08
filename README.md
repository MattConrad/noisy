# noisy
Simple server in Go that plays wavfiles and accepts short inputs for AT&amp;T speech API.

This project allows you to control a remote computer and have it play sound files. I use it to talk to my family when I'm downstairs and they're upstairs.

These sound files can be created ahead of time, or created on the fly using the AT&T speech API.

The play() function presently only supports Windows machines that have PowerShell installed. If you have a Win 7+ machine, it probably has PowerShell installed already.

If you have predefined .wav files you call them by POSTing to an URL /run/:wavname. See routing in main() and the run() function: POSTing to /run/cindy_telephone plays the .wav file cindy_telephone.wav. The directory for prerecorded .wav files is presently hardcoded to ./wavfiles/.

You can also enter arbitrary text to be converted into a .wav by the AT&T speech API. To use the speech API, you'll need an app key id and app secret, which are presently passed into the app as flags. Example:

  noisy -attid=v6fcaxkskdfszzz23qrgb8aob3uy2 -attsecret=f4zkec5uxcvzzzkst91hnyq73qrbz

