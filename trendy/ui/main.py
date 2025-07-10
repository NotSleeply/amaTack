from datetime import datetime
import random
import tkinter as tk
from tkinter import messagebox, filedialog
import json, threading, websocket
import time
import os

# ────────── 常量 ──────────
NOGO_NAME = "Trendy"
BOARD_SIZE = 9
CELL_SIZE = 48
MARGIN = CELL_SIZE
STONE_R = CELL_SIZE // 2 - 3
WS_ADDR = "ws://localhost:8080/ws"

BLACK, WHITE, EMPTY = 1, -1, 0

# ────────── GUI ──────────
class NoGoGUI(tk.Tk):
    def __init__(self):
        super().__init__()
        self.title(NOGO_NAME)
        self.resizable(False, False)
        self.configure(bg="#F8DCA7")

        # 游戏状态
        self.board = [[EMPTY]*BOARD_SIZE for _ in range(BOARD_SIZE)]
        self.current = BLACK
        self.ai_color = WHITE
        self.game_over = False
        self.ai_connected = False
        self.move_history = []  # 历史记录: [(row, col, color), ...]
        self.move_count = 0     # 手数计数
        self.ai_think_time = "3" # 默认3秒
        # 为每局游戏生成唯一的棋谱文件名
        self.filename = self._generate_filename()

        # 计时器 - 修复版本
        self.black_time = 15 * 60  # 15分钟 = 900秒
        self.white_time = 15 * 60
        self.last_update_time = None  # 上次更新时间
        self.timer_running = False

        # 画布
        size = MARGIN*2 + (BOARD_SIZE-1)*CELL_SIZE
        self.canvas = tk.Canvas(self, width=size, height=size,
                                bg="#F8DCA7", highlightthickness=0)
        self.canvas.pack()
        self._draw_board()
        self.canvas.bind("<Button-1>", self.on_click)

        # 信息标签
        self.info_lbl = tk.Label(self, font=("Arial", 12), bg="#F8DCA7")
        self.info_lbl.pack(pady=4)
        self._update_info()

        # 菜单
        menubar = tk.Menu(self)
        menubar.add_command(label="新游戏(AI后手)", command=lambda: self._new_game(WHITE))
        menubar.add_command(label="新游戏(AI先手)", command=lambda: self._new_game(BLACK))
        menubar.add_command(label="人人对战", command=lambda: self._new_game(None))
        menubar.add_command(label="悔棋", command=self._undo)
        menubar.add_command(label="保存棋谱", command=self._manual_save_records)
        # 添加AI思考时间菜单
        ai_time_menu = tk.Menu(menubar, tearoff=0)
        ai_time_menu.add_command(label="1秒", command=lambda: self._set_ai_think_time("1"))
        ai_time_menu.add_command(label="3秒", command=lambda: self._set_ai_think_time("3"))
        ai_time_menu.add_command(label="5秒", command=lambda: self._set_ai_think_time("5"))
        ai_time_menu.add_command(label="10秒", command=lambda: self._set_ai_think_time("10"))
        ai_time_menu.add_command(label="15秒", command=lambda: self._set_ai_think_time("15"))
        ai_time_menu.add_command(label="30秒", command=lambda: self._set_ai_think_time("30"))
        menubar.add_cascade(label="AI思考时间", menu=ai_time_menu)
        self.config(menu=menubar)

        # WebSocket
        self.ws = None
        threading.Thread(target=self._ws_connect, daemon=True).start()

        # 启动计时器
        self._start_timer()

    def _generate_filename(self):
        """生成唯一的棋谱文件名"""
        timestamp = datetime.now().strftime('%H时%M分%S秒')
        return f"先手队名 vs 后手队名-先后手胜_{timestamp}.txt"

    def _set_ai_think_time(self, seconds):
        """设置AI思考时间"""
        self.ai_think_time = seconds

    def _ws_connect(self):
        while True:
            try:
                self.ws = websocket.WebSocket()
                self.ws.connect(WS_ADDR)
                self.ai_connected = True
                self._update_info()

                if self.current == self.ai_color:
                    self._send_board()

                while True:
                    msg = self.ws.recv()
                    data = json.loads(msg)
                    if isinstance(data, list) and len(data) == 2:
                        row, col = data
                        if row== -1 and col == -1:
                           self._make_random_move()
                        else:
                            self.after(0, self._ai_move, row, col)

            except Exception as e:
                print(f"WebSocket错误: {e}")
                self.ai_connected = False
                self._update_info()
                time.sleep(2)

    def _send_board(self):
        # 只有在AI连接且轮到AI时才发送棋盘
        if self.ws and self.ai_connected and self.current == self.ai_color and not self.game_over:
            data = {"board": self.board, "time": self.ai_think_time}
            self.ws.send(json.dumps(data))

    def _draw_board(self):
        self.canvas.delete("all")

        # 画线
        for i in range(BOARD_SIZE):
            p = MARGIN + i*CELL_SIZE
            self.canvas.create_line(p, MARGIN, p, MARGIN+(BOARD_SIZE-1)*CELL_SIZE)
            self.canvas.create_line(MARGIN, p, MARGIN+(BOARD_SIZE-1)*CELL_SIZE, p)

        # 画星位
        for r,c in [(4,4),(2,2),(2,6),(6,2),(6,6)]:
            x = MARGIN + c*CELL_SIZE
            y = MARGIN + r*CELL_SIZE
            self.canvas.create_oval(x-3, y-3, x+3, y+3, fill="black")

        # 画棋子和序号
        for r in range(BOARD_SIZE):
            for c in range(BOARD_SIZE):
                if self.board[r][c] != EMPTY:
                    move_num = self._get_move_number(r, c)
                    self._draw_stone(r, c, self.board[r][c], move_num)

    def _get_move_number(self, r, c):
        """获取指定位置棋子的手数"""
        for i, (hr, hc, _) in enumerate(self.move_history):
            if hr == r and hc == c:
                return i + 1
        return 0

    def _draw_stone(self, r, c, color, move_num=0):
        x = MARGIN + c*CELL_SIZE
        y = MARGIN + r*CELL_SIZE
        fill = "black" if color == BLACK else "white"
        outline = "white" if color == BLACK else "black"

        # 画棋子
        self.canvas.create_oval(x-STONE_R, y-STONE_R, x+STONE_R, y+STONE_R,
                                fill=fill, outline=outline, width=2)

        # 画序号
        if move_num > 0:
            # 判断是否为最后一步棋
            is_last_move = (move_num == len(self.move_history))

            if is_last_move:
                # 最后一步棋的数字用红色显示
                text_color = "red"
            else:
                # 其他棋子的数字用常规颜色
                text_color = "white" if color == BLACK else "black"

            self.canvas.create_text(x, y, text=str(move_num),
                                   fill=text_color, font=("Arial", 10, "bold"))

    def on_click(self, event):
        if self.game_over:
            return

        # AI未连接时允许人人对战，AI连接时只能在非AI回合下棋
        if self.ai_connected and self.current == self.ai_color:
            return

        c = round((event.x - MARGIN) / CELL_SIZE)
        r = round((event.y - MARGIN) / CELL_SIZE)

        if 0 <= r < BOARD_SIZE and 0 <= c < BOARD_SIZE and self.board[r][c] == EMPTY:
            self._make_move(r, c)

    def _ai_move(self, r, c):
        # 只有在AI模式、轮到AI、且AI已连接时才处理AI移动
        if (not self.game_over and self.ai_color is not None and
            self.current == self.ai_color and self.ai_connected and
            self.board[r][c] == EMPTY):
            self._make_move(r, c)

    def _make_move(self, r, c):
        # 停止计时并记录用时
        if self.timer_running:
            self._stop_timer()

        # 记录历史
        self.move_history.append((r, c, self.current))
        self.move_count += 1

        # 下棋
        self.board[r][c] = self.current

        # 重新绘制整个棋盘以更新数字颜色
        self._draw_board()

        # 检查游戏是否结束
        if self._check_game_over():
            self.game_over = True
            winner = "黑方" if self.current == WHITE else "白方"
            # 自动保存棋谱
            self._save_records()
            messagebox.showinfo("游戏结束", f"{winner}获胜！\n棋谱已保存至: {self.filename}")
            return

        # 切换回合
        self.current = -self.current
        self._start_timer()  # 重新开始计时
        self._update_info()

        # 只有AI模式才发送棋盘给AI
        if self.ai_color is not None:
            self._send_board()

    def _make_random_move(self):
        """AI认输后的随机下棋"""
        if self.game_over:
            return

        # 获取所有空位
        empty_positions = []
        for r in range(BOARD_SIZE):
            for c in range(BOARD_SIZE):
                if self.board[r][c] == EMPTY:
                    empty_positions.append((r, c))

        if empty_positions:
            # 随机选择一个空位
            r, c = random.choice(empty_positions)
            # 延迟一下让用户看到是随机下的
            self.after(500, lambda: self._make_move(r, c))

    def _undo(self):
        """悔棋功能"""
        if not self.move_history or self.game_over:
            return

        # 停止计时
        if self.timer_running:
            self._stop_timer()

        if self.move_history:
            r, c, _ = self.move_history.pop()
            self.board[r][c] = EMPTY
            self.move_count -= 1
            self.current = -self.current

        self._draw_board()
        self._start_timer()  # 重新开始计时
        self._update_info()

    def _check_game_over(self):
        """检查游戏是否结束（当前玩家是否违规）"""
        for r in range(BOARD_SIZE):
            for c in range(BOARD_SIZE):
                if self.board[r][c] == self.current:
                    if not self._has_liberty(r, c, set()):
                        return True
                    for dr, dc in [(0,1), (1,0), (0,-1), (-1,0)]:
                        nr, nc = r + dr, c + dc
                        if (0 <= nr < BOARD_SIZE and 0 <= nc < BOARD_SIZE and
                            self.board[nr][nc] == -self.current):
                            if not self._has_liberty(nr, nc, set()):
                                return True
        self._save_records()
        return False

    def _has_liberty(self, r, c, visited):
        """检查棋子或棋组是否有气"""
        if (r, c) in visited:
            return False
        visited.add((r, c))

        color = self.board[r][c]
        for dr, dc in [(0,1), (1,0), (0,-1), (-1,0)]:
            nr, nc = r + dr, c + dc
            if 0 <= nr < BOARD_SIZE and 0 <= nc < BOARD_SIZE:
                if self.board[nr][nc] == EMPTY:
                    return True
                elif self.board[nr][nc] == color:
                    if self._has_liberty(nr, nc, visited):
                        return True
        return False

    def _new_game(self, ai_color):
        self.board = [[EMPTY]*BOARD_SIZE for _ in range(BOARD_SIZE)]
        self.current = BLACK
        self.ai_color = ai_color  # None表示人人对战模式
        self.game_over = False
        self.move_history = []
        self.move_count = 0
        self.black_time = 15 * 60
        self.white_time = 15 * 60
        # 为新游戏生成新的文件名
        self.filename = self._generate_filename()
        self._draw_board()
        self._start_timer()
        self._update_info()

        # 只有在AI模式且AI先手时才发送棋盘
        if self.ai_color == BLACK and self.ai_connected:
            self._send_board()

    def _update_info(self):
        current_player = "黑方" if self.current == BLACK else "白方"

        # 根据连接状态显示不同信息
        if self.ai_connected:
            ai_player = "黑方" if self.ai_color == BLACK else "白方"
            mode_info = f"AI执{ai_player}| 思考：{self.ai_think_time}秒 "
        else:
            mode_info = "人人对战 "

        # 格式化时间显示
        black_min, black_sec = divmod(int(self.black_time), 60)
        white_min, white_sec = divmod(int(self.white_time), 60)

        self.info_lbl.config(text=f"当前：{current_player} | {mode_info} | 手数：{self.move_count}\n"
                                  f"黑方时间：{black_min:02d}:{black_sec:02d} | 白方时间：{white_min:02d}:{white_sec:02d}")

    def _start_timer(self):
        """开始计时"""
        if not self.game_over:
            self.last_update_time = time.time()
            self.timer_running = True
            self._update_timer()

    def _stop_timer(self):
        """停止计时并更新剩余时间"""
        if self.timer_running and self.last_update_time:
            elapsed = time.time() - self.last_update_time
            if self.current == BLACK:
                self.black_time = max(0, self.black_time - elapsed)
            else:
                self.white_time = max(0, self.white_time - elapsed)
            self.timer_running = False

    def _update_timer(self):
        """更新计时器显示 - 修复版本"""
        if self.timer_running and not self.game_over:
            current_time = time.time()

            # 计算这次更新的时间间隔（应该约为1秒）
            if self.last_update_time:
                elapsed = current_time - self.last_update_time

                # 从对应玩家的时间中减去这个时间间隔
                if self.current == BLACK:
                    self.black_time = max(0, self.black_time - elapsed)
                    if self.black_time <= 0:
                        # self.game_over = True
                        # 自动保存棋谱
                        self._save_records()
                        messagebox.showinfo("游戏结束", f"黑方超时，白方获胜！\n棋谱已保存至: {self.filename}")
                        return
                else:
                    self.white_time = max(0, self.white_time - elapsed)
                    if self.white_time <= 0:
                        # self.game_over = True
                        # 自动保存棋谱
                        self._save_records()
                        messagebox.showinfo("游戏结束", f"白方超时，黑方获胜！\n棋谱已保存至: {self.filename}")
                        return

            # 更新"上次更新时间"为当前时间
            self.last_update_time = current_time

            # 更新显示
            self._update_info()

            # 1秒后再次调用
            self.after(1000, self._update_timer)

    def _save_records(self, record_file=None):
        """保存当前棋谱到 TXT 文件"""
        if record_file is None:
            record_file = self.filename

        def convert_to_chess_notation(history):
            board_size = BOARD_SIZE
            moves = []
            for i, (row, col, color) in enumerate(history):
                player = 'B' if color == BLACK else 'W'
                col_letter = chr(ord('A') + col)
                row_number = board_size - row
                moves.append(f"{player}[{col_letter}{row_number}]")
            return ";".join(moves)

        try:
            if not self.move_history:
                print("没有棋谱数据可保存")
                return False

            moves_str = convert_to_chess_notation(self.move_history)
            now_str = datetime.now().strftime("%Y-%m-%d %H:%M")

            # 确定游戏结果
            if self.game_over:
                if self.black_time <= 0:
                    result = "白胜"
                elif self.white_time <= 0:
                    result = "黑胜"
                else:
                    # 正常游戏结束，当前玩家违规
                    result = "黑胜" if self.current == WHITE else "白胜"
            else:
                result = "未结束"

            # 使用正确的棋谱格式：方括号包围整个棋谱，内部信息也用方括号
            record_str = f"([NG][Trendy][AI][{result}][{now_str}][CCGC];{moves_str})"

            with open(record_file, "w", encoding="utf-8") as f:
                f.write(record_str)

            print(f"棋谱保存至 {record_file} 成功")
            return True

        except Exception as e:
            print(f"保存棋谱失败: {e}")
            return False

    def _manual_save_records(self):
        """手动保存棋谱，允许用户选择文件名和位置"""
        if not self.move_history:
            messagebox.showwarning("保存失败", "没有棋谱数据可保存")
            return

        # 打开文件保存对话框
        file_path = filedialog.asksaveasfilename(
            title="保存棋谱",
            defaultextension=".txt",
            filetypes=[("文本文件", "*.txt"), ("所有文件", "*.*")],
            initialfile=self.filename
        )

        if file_path:
            if self._save_records(file_path):
                messagebox.showinfo("保存成功", f"棋谱已保存至:\n{file_path}")
            else:
                messagebox.showerror("保存失败", "保存棋谱时发生错误")

if __name__ == "__main__":
    NoGoGUI().mainloop()
