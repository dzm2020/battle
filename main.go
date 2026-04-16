package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"battle/internal/battle/attr"
	"battle/internal/battle/calc"
	"battle/internal/battle/clock"
	"battle/internal/battle/control"
	"battle/internal/battle/cooldown"
	"battle/internal/battle/entity"
	"battle/internal/battle/geom"
	"battle/internal/battle/room"
	"battle/internal/battle/skill"
	"battle/internal/battle/tick"
	"battle/internal/battle/timer"
)

// -------------------------- 【1】服务端配置 --------------------------
const (
	HeaderSize  = 4       // 消息头：4字节表示消息长度
	ReadTimeout = 10 * time.Second
)

// -------------------------- 【2】客户端连接实体 --------------------------
// Client 代表一个客户端连接（玩家/客户端）
type Client struct {
	Conn     net.Conn
	ID       string
	LastTime int64 // 心跳时间
}

// -------------------------- 【3】全局管理器 --------------------------
var (
	clientMap = make(map[string]*Client)
	clientMu  sync.RWMutex
)

func main() {
	fmt.Println("=== Go 战斗服务器 第1～5天：网络 + 属性 + 帧循环 + 房间 + 技能（mock 客户端模式）===")

	day02Demo()
	day03Demo()
	day04Demo()
	day05Demo()

	serverConn, clientConn := net.Pipe()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\n服务器关闭中...")
		closeAllClient()
		_ = clientConn.Close()
		_ = serverConn.Close()
		os.Exit(0)
	}()

	clientID := "mock_client"
	client := &Client{
		Conn:     serverConn,
		ID:       clientID,
		LastTime: time.Now().Unix(),
	}

	clientMu.Lock()
	clientMap[clientID] = client
	clientMu.Unlock()

	fmt.Printf("mock 客户端已接入：%s\n", clientID)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		handleClient(client)
	}()
	go func() {
		defer wg.Done()
		runMockClient(clientConn)
	}()
	wg.Wait()
	fmt.Println("mock 会话结束，进程退出")
}

func day02Demo() {
	c := calc.DefaultCalculator{}
	hero := entity.New("hero_1", 1, attr.Base{
		Level: 10,
		STR:   20, AGI: 15, INT: 12, VIT: 18,
	})
	hero.InitBattle(c)

	fmt.Println("--- 第2天：Entity / Base / Derived / Runtime ---")
	fmt.Printf("实体 %s Camp=%d\n", hero.ID, hero.Camp)
	fmt.Printf("Base: %+v\n", hero.Base)
	fmt.Printf("Derived: HP=%d/%d MP=%d/%d ATK=%d DEF=%d Crit=%.3f CritDmg=%.2f PhysMit=%.3f\n",
		hero.Runtime.CurHP, hero.Derived.MaxHP,
		hero.Runtime.CurMP, hero.Derived.MaxMP,
		hero.Derived.ATK, hero.Derived.DEF,
		hero.Derived.CritRate, hero.Derived.CritDamage, hero.Derived.PhysMitigation,
	)

	hero.Runtime.CurHP = 30
	hero.Base.VIT = 30
	hero.Recalculate(c)
	fmt.Printf("升体质后裁剪血量: CurHP=%d MaxHP=%d\n", hero.Runtime.CurHP, hero.Derived.MaxHP)
	fmt.Println()
}

func day03Demo() {
	fmt.Println("--- 第3天：Clock / Tick Loop / Timer / Cooldown ---")
	clk := clock.New(60)
	loop := tick.NewLoop(clk)
	tm := timer.NewManager()
	cd := cooldown.NewBook()

	const (
		tagDelayed timer.Tag = 1
		tagPulse   timer.Tag = 2
	)
	tm.AddOneShot(30, tagDelayed)
	tm.AddRepeat(20, 15, tagPulse)

	loop.Add(tick.FuncSubscriber(func(c *clock.Clock) {
		for _, ev := range tm.ProcessFrame(c.Frame()) {
			fmt.Printf("  [frame=%d ms=%d] timer id=%d tag=%d\n", c.Frame(), c.LogicalMs(), ev.ID, ev.Tag)
		}
	}))

	for i := 0; i < 45; i++ {
		if i == 0 {
			fmt.Printf("  skill fireball ready before cast: %v\n", cd.IsReady(clk.Frame(), "fireball"))
		}
		loop.Step()
		if clk.Frame() == 10 {
			cd.Trigger(clk.Frame(), "fireball", 25)
			fmt.Printf("  [frame=%d] cast fireball, cd=25f\n", clk.Frame())
		}
		if clk.Frame() == 12 || clk.Frame() == 35 {
			nextStr := "n/a"
			if n, ok := cd.NextReadyFrame("fireball"); ok {
				nextStr = fmt.Sprintf("%d", n)
			}
			fmt.Printf("  [frame=%d] fireball ready=%v next=%s\n", clk.Frame(),
				cd.IsReady(clk.Frame(), "fireball"), nextStr)
		}
	}
	fmt.Printf("  共推进 %d 帧，逻辑时长约 %d ms\n", clk.Frame(), clk.LogicalMs())
	fmt.Println()
}

func day04Demo() {
	fmt.Println("--- 第4天：BattleManager / Room 生命周期 ---")
	mgr := room.NewManager()
	r, err := mgr.Create("room_alpha", 2)
	if err != nil {
		fmt.Println("  Create:", err)
		return
	}
	cal := calc.DefaultCalculator{}
	p1 := entity.New("hero_p1", 1, attr.Base{Level: 5, STR: 10, AGI: 8, INT: 6, VIT: 12})
	p2 := entity.New("hero_p2", 2, attr.Base{Level: 5, STR: 8, AGI: 10, INT: 8, VIT: 10})
	_ = r.Join("sess_1", p1)
	_ = r.Join("sess_2", p2)
	fmt.Printf("  room=%s phase=%s players=%d/%d\n", r.ID(), r.Phase(), r.PlayerCount(), r.MaxPlayers())

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	if err := r.StartBattle(ctx, cal); err != nil {
		fmt.Println("  StartBattle:", err)
		return
	}
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("  after tick phase=%s\n", r.Phase())

	if err := r.Settle(); err != nil {
		fmt.Println("  Settle:", err)
	} else {
		fmt.Printf("  Settle ok phase=%s\n", r.Phase())
	}
	mgr.Destroy("room_alpha")
	fmt.Printf("  manager room count=%d\n", mgr.Count())
	fmt.Println()
}

func day05Demo() {
	fmt.Println("--- 第5天：技能系统（校验链 / 沉默 / 前摇 / 自身目标）---")
	clk := clock.New(60)
	loop := tick.NewLoop(clk)
	reg := skill.MustDemoRegistry()
	sys := skill.NewSystem(reg, skill.DefaultApplier{})
	loop.Add(sys)

	cal := calc.DefaultCalculator{}
	caster := entity.New("caster", 1, attr.Base{Level: 10, STR: 15, AGI: 10, INT: 20, VIT: 12})
	foe := entity.New("foe", 2, attr.Base{Level: 10, STR: 10, AGI: 10, INT: 10, VIT: 10})
	for _, id := range []string{"strike", "fireball", "focus"} {
		caster.GrantSkill(id)
	}
	caster.InitBattle(cal)
	foe.InitBattle(cal)
	caster.Pos = geom.Vec2{X: 0, Y: 0}
	foe.Pos = geom.Vec2{X: 2, Y: 0}

	loop.Step()
	res := sys.TryCast(skill.CastInput{Frame: clk.Frame(), BattleActive: true, Caster: caster, Target: foe, SkillID: "strike"})
	fmt.Printf("  strike: ok=%v stage=%v reason=%v mp=%d\n", res.OK, res.Stage, res.Reason, caster.Runtime.CurMP)

	caster.Control = control.FlagSilenced
	res2 := sys.TryCast(skill.CastInput{Frame: clk.Frame(), BattleActive: true, Caster: caster, Target: foe, SkillID: "fireball"})
	fmt.Printf("  fireball silenced: ok=%v reason=%v\n", res2.OK, res2.Reason)
	caster.Control = 0

	res3 := sys.TryCast(skill.CastInput{Frame: clk.Frame(), BattleActive: true, Caster: caster, Target: foe, SkillID: "fireball"})
	fmt.Printf("  fireball start: ok=%v stage=%v endsFrame=%d mp=%d\n", res3.OK, res3.Stage, res3.WindupEndsAtFrame, caster.Runtime.CurMP)
	for clk.Frame() < res3.WindupEndsAtFrame {
		loop.Step()
	}
	fmt.Printf("  after windup mp=%d foe_hp=%d\n", caster.Runtime.CurMP, foe.Runtime.CurHP)

	res4 := sys.TryCast(skill.CastInput{Frame: clk.Frame(), BattleActive: true, Caster: caster, Target: nil, SkillID: "focus"})
	fmt.Printf("  focus self: ok=%v reason=%v mp=%d\n", res4.OK, res4.Reason, caster.Runtime.CurMP)
	fmt.Println()
}

// runMockClient 在 Pipe 另一端模拟客户端，按与服务端相同的帧格式收发
func runMockClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	steps := []struct {
		send string
		want string
	}{
		{"ping", "pong"},
		{"enter_battle", "battle_enter_ok"},
		{"use_skill", "skill_not_implemented"},
		{"nope", "unknown_msg"},
	}

	for _, step := range steps {
		if err := writeFramed(conn, step.send); err != nil {
			fmt.Println("[mock] 发送失败：", err)
			return
		}
		fmt.Printf("[mock] 发送：%s\n", step.send)

		reply, err := readFramed(reader)
		if err != nil {
			fmt.Println("[mock] 收包失败：", err)
			return
		}
		fmt.Printf("[mock] 收到：%s", reply)
		if reply != step.want {
			fmt.Printf("（期望 %s）", step.want)
		}
		fmt.Println()
		time.Sleep(50 * time.Millisecond)
	}
}

func writeFramed(w io.Writer, msg string) error {
	msgBytes := []byte(msg)
	msgLen := uint32(len(msgBytes))
	buf := make([]byte, HeaderSize+msgLen)
	binary.BigEndian.PutUint32(buf[:HeaderSize], msgLen)
	copy(buf[HeaderSize:], msgBytes)
	_, err := w.Write(buf)
	return err
}

func readFramed(r *bufio.Reader) (string, error) {
	headerBuf := make([]byte, HeaderSize)
	if _, err := io.ReadFull(r, headerBuf); err != nil {
		return "", err
	}
	msgLen := binary.BigEndian.Uint32(headerBuf)
	if msgLen <= 0 || msgLen > 1024*1024 {
		return "", fmt.Errorf("消息长度非法: %d", msgLen)
	}
	msgBody := make([]byte, msgLen)
	if _, err := io.ReadFull(r, msgBody); err != nil {
		return "", err
	}
	return string(msgBody), nil
}

// -------------------------- 【4】客户端消息处理（核心） --------------------------
func handleClient(client *Client) {
	defer func() {
		client.Conn.Close()
		clientMu.Lock()
		delete(clientMap, client.ID)
		clientMu.Unlock()
		fmt.Printf("客户端断开：%s | 剩余在线：%d\n", client.ID, len(clientMap))
	}()

	reader := bufio.NewReader(client.Conn)

	for {
		client.Conn.SetReadDeadline(time.Now().Add(ReadTimeout))

		msgStr, err := readFramed(reader)
		if err != nil {
			if err != io.EOF {
				fmt.Println("读包失败：", err)
			}
			return
		}

		client.LastTime = time.Now().Unix()
		fmt.Printf("[%s] 收到消息：%s\n", client.ID, msgStr)
		dispatchMsg(client, msgStr)
	}
}

// -------------------------- 【5】消息分发（战斗协议入口） --------------------------
func dispatchMsg(client *Client, msg string) {
	switch msg {
	case "ping":
		sendMsg(client, "pong")
	case "enter_battle":
		sendMsg(client, "battle_enter_ok")
	case "use_skill":
		sendMsg(client, "skill_not_implemented")
	default:
		sendMsg(client, "unknown_msg")
	}
}

// -------------------------- 【6】发送消息 --------------------------
func sendMsg(client *Client, msg string) {
	_ = writeFramed(client.Conn, msg)
}

// 关闭所有客户端
func closeAllClient() {
	clientMu.RLock()
	defer clientMu.RUnlock()

	for _, c := range clientMap {
		c.Conn.Close()
	}
}
