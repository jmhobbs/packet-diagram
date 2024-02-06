# Packet Diagram

This is intended to be a vizualizer for binary data, communication packets specifically.  It is very simple, and has a simple format.

## Config

Describing a packet is done in a simple text file. Each line is a rule and a description.

Currently, there are four rules:

| `N` | Return the next `N` bytes |
| `-N` | Skip the next `N` bytes |
| `>XX` | Return until the next instance of the byte value `XX` (in base 16) |
| `->XX` | Skip until the next instance of the byte value `XX` (in base 16) |

The input file is then sliced up based on this description and printed with offset, bytes, printables and description.

Bytes remaining after the end of the description are discarded.

## Example

### Packet Description File

```
4 Response Type
1 Packet Type
1 Protocol
>0 Name
>0 Map
```

### Binary File (hex dumped for printing here)

```
0000000      ffff    ffff    1149    6144    5a79    5520    2053    202d
0000020      594e    3620    3530    2033    3128    7473    5020    7265
0000040      6f73    206e    6e4f    796c    0029    6863    7265    616e
0000060      7572    7073    756c    0073    6164    7a79    4400    7961
0000100      005a    0000    3c23    6400    0077    3101    322e    2e33
0000120      3531    3037    3534    b100    2774    3c03    93ad    62b4
0000140      0140    6162    7474    656c    6579    6e2c    336f    6472
0000160      732c    6168    6472    3030    2c31    716c    3073    652c
0000200      6d74    2e34    3032    3030    3030    652c    746e    346d
0000220      302e    3030    3030    2c30    3431    303a    0039    5fac
0000240      0003    0000    0000
0000246
```

### Output

```
┏━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━┯━━━━━━━━┯━━━━━━━━━━━━━━━┓
┃00000000│ ff ff ff ff             ┊                         │⋄⋄⋄⋄    ┊        │ Response Type ┃
┠┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┨
┃00000004│ 49                      ┊                         │I       ┊        │ Packet Type   ┃
┠┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┨
┃00000005│ 11                      ┊                         │⋄       ┊        │ Protocol      ┃
┠┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┨
┃00000006│ 44 61 79 5a 20 55 53 20 ┊ 2d 20 4e 59 20 36 30 35 │DayZ US ┊- NY 605│ Name          ┃
┃00000022│ 33 20 28 31 73 74 20 50 ┊ 65 72 73 6f 6e 20 4f 6e │3 (1st P┊erson On│               ┃
┃00000038│ 6c 79 29 00             ┊                         │ly)⋄    ┊        │               ┃
┠┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┨
┃00000042│ 63 68 65 72 6e 61 72 75 ┊ 73 70 6c 75 73 00       │chernaru┊splus⋄  │ Map           ┃
┗━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━┷━━━━━━━━┷━━━━━━━━━━━━━━━┛
```

## TODO

  - [ ] Implement skipping in presentation
  - [ ] Colorization
  - [ ] Additional rules
  - [ ] Split parsing from presentation
  - [ ] Aternative renderers (HTML, Markdown, ???)
  - [ ] Stream packet reading
