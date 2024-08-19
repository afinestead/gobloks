import math

matrix = [
  [0, 0, 0, 0, 0, 0, 0, 0],
  [0]*8,
  [0]*8,
  [0]*8,
  [0]*8,
  [0, 0, 1, 0, 0, 0, 0, 0],
  [1, 1, 1, 0, 0, 0, 0, 0],
  [0, 1, 0, 0, 0, 0, 0, 0]
]
# matrix = [
#   [1, 0, 0, 0, 0, 0, 0, 0],
#   [0]*8,
#   [0]*8,
#   [0]*8,
#   [0]*8,
#   [0]*8,
#   [0]*8,
#   [0]*8,
# ]

to_int = lambda mtx: int("".join([str(element) for sublist in reversed(mtx) for element in reversed(sublist)]), 2)
def from_int(num):
  as_bin = format(num, '064b')
  return [[int(bit) for bit in reversed(as_bin[i:i+8])] for i in range(0, 64, 8)[::-1]]

def pretty_print(mtx):
  s = ""
  for row in mtx:
    s += " ".join([str(bit) for bit in row]) + "\n"
  print(s)


MAX_DEGREE = 8
MAGIC = 0x02040810204081
COLUMN_MASK = 0x8080808080808080

def get_column(n: int, col: int) -> int:
  return ((((n<<(MAX_DEGREE-1-col))&COLUMN_MASK)*MAGIC)>>((MAX_DEGREE**2)-MAX_DEGREE))&((2**MAX_DEGREE)-1)

def _clz(n: int):
  idx = 0
  while not (n & 1):
    n >>= 1
    idx += 1
  return idx

def rotate90AndNormalize(num):
  out = 0
  col_mask = 0x0101010101010101
  lsx = lsy = MAX_DEGREE**2
  for i in range(8):
    col = num & (col_mask << i)
    row = 0
    for j in range(8):
      bit = (col >> ((i + (j*8)) - j)) & 0xff
      if bit:
        lsx = min(lsx, j)

      row |= bit
    row_pos = ((8-i-1)*8)
    if row:
      lsy = min(lsy, row_pos)
    out |= row << row_pos
  return out >> (lsx+lsy)

def rotate90AndNormalize2(num):
  out = 0
  lsx = lsy = MAX_DEGREE**2
  for c in range(MAX_DEGREE):
    col = get_column(num, c)
    row_shift = ((MAX_DEGREE - c - 1) * MAX_DEGREE)
    if col > 0:
      lsx = min(lsx, _clz(col))
      lsy = min(lsy, row_shift)
    out |= (col << row_shift)
  return out >> (lsx+lsy)


# as_int = to_int(matrix)
# print(as_int)
# # # assert(as_int == 1)
# # # print(bin(as_int))
# # # print(from_int(as_int))
# # # assert(from_int(as_int) == matrix)
# # pretty_print(matrix)
# pretty_print(from_int(as_int))

# # # print(bin(res))
# # res = rotate90AndNormalize(as_int)
# # pretty_print(from_int(res))
# # res2 = rotate90AndNormalize(res)
# # pretty_print(from_int(res2))
# # res3 = rotate90AndNormalize(res2)
# # pretty_print(from_int(res3))
# # res4 = rotate90AndNormalize(res3)
# # pretty_print(from_int(res4))

# res = rotate90AndNormalize2(as_int)
# print(bin(res))
# pretty_print(from_int(res))
# res2 = rotate90AndNormalize2(res)
# print(bin(res2))
# pretty_print(from_int(res2))
# res3 = rotate90AndNormalize2(res2)
# print(bin(res3))
# pretty_print(from_int(res3))
# res4 = rotate90AndNormalize2(res3)
# print(bin(res4))
# pretty_print(from_int(res4))


# for ii in range(2**64):
#   print(ii)
#   if to_col(col_mask*ii) == 0xff:
#     assert(hex(ii))

import matplotlib.pyplot as plt
points = [
  (10, 20),
  (0, 11),
  (8, 20),
  (20, 7),
  (14, 1),
  (11, 20),
  (0, 9),
  (14, 19),
  (6, 19),
  (18, 16),
  (3, 17),
  (12, 20),
  (0, 13),
  (0, 12),
  (11, 0),
  (9, 0),
  (6, 1),
  (19, 6),
  (4, 2),
  (20, 10),
  (8, 0),
  (20, 8),
  (20, 9),
  (7, 20),
  (1, 6),
  (15, 19),
  (2, 16),
  (3, 3),
  (1, 15),
  (17, 3),
  (0, 10),
  (19, 15),
  (19, 5),
  (1, 5),
  (2, 4),
  (13, 0),
  (5, 19),
  (9, 20),
  (20, 11),
  (0, 7),
  (4, 18),
  (16, 2),
  (18, 4),
  (17, 17),
  (20, 12),
  (0, 8),
  (13, 20),
  (15, 1),
  (10, 0),
  (12, 0),
  (7, 0),
  (19, 14),
  (1, 14),
  (5, 1),
  (20, 13),
  (16, 18),
]

x = [pt[0] for pt in points]
y = [pt[1] for pt in points]

# circle.center.X + int(circle.radius)*int(math.Cos(theta)),
# circle.center.Y + int(circle.radius)*int(math.Sin(theta)),
# theta = math.pi / 4
# # theta = 
# pt = (
#   int(10*math.cos(theta)),
#   int(10*math.sin(theta)),
# )

# print(pt)

plt.scatter(x, y)
plt.show()