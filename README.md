# World Trimmer
This is a simple project to trim a world by deleting region files outside of a certain amount of specified blocks. It operates by deleting .mca files in the region folder of your world by calculating the bounds of the region files.

This project would be useful for individuals who have massive world files due to not using a plugin such as worldborder, but would like to keep the same world and delete the excess regions.

## Usage
Download the binary from this github repository and run the following command:

```
./wc threshold debug path
./wc 25000 false /root/myserver/world/region
```
Make sure to specify the path exactly into the folder where the .mca files are located.
If debug is true it will not ACTUALLY delete the files. It will show you how many files it would delete, and how much storage it would save you.

## Notes
Minecraft chunk data storage is stored in .mca files. Each .mca file stores 512x512 blocks. That means that there is no way to trim a world exactly to a nice number such as 10,000. My code will ensure that no region file will be deleted before the threshold provided.

Most of my math was checked by https://dinnerbone.com/minecraft/tools/coordinates/. Thanks to dinnerbone for creating that lovely tool.
I would also like to thank Bandwidth from SkittleMC for allowing me to open source this project as this was originally just for him.

This was also my first attempt at a golang project. Constructive criticism would be much appreciated.

If anyone needs any assistance using this project I would love to help. Just PM me on twitter or Spigot. @RenegadeEagle
