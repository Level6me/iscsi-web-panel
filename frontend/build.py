#!/usr/bin/env python3
"""Generate complete iSCSI Web Panel index.html"""

html = '''<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1,user-scalable=no,viewport-fit=cover">
<meta name="theme-color" content="#f2f2f7" media="(prefers-color-scheme:light)">
<meta name="theme-color" content="#000" media="(prefers-color-scheme:dark)">
<title>iSCSI Panel</title>
<script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
<style>
:root{--bg:#f2f2f7;--card:#fff;--text:#000;--sec:#8e8e93;--accent:#007aff;--danger:#ff3b30;--ok:#34c759;--warn:#ff9500;--purple:#af52de;--dock:rgba(255,255,255,.85);--brd:rgba(128,128,128,.15)}
@media(prefers-color-scheme:dark){:root{--bg:#000;--card:#1c1c1e;--text:#fff;--sec:#98989d;--dock:rgba(28,28,30,.85)}}
*{box-sizing:border-box;margin:0;padding:0}
body{background:var(--bg);color:var(--text);font-family:-apple-system,BlinkMacSystemFont,'SF Pro',sans-serif;padding-bottom:calc(100px + env(safe-area-inset-bottom))}
.ctn{padding:20px 16px;max-width:600px;margin:0 auto;padding-top:max(20px,env(safe-area-inset-top))}
.hdr{margin:10px 0 20px;display:flex;justify-content:space-between;align-items:flex-end}
.ttl{font-size:34px;font-weight:700}
.hdr-sub{font-size:12px;color:var(--sec)}
.hdr-r{text-align:right}
.sdot{font-size:12px;color:var(--ok)}
.utime{font-size:11px;opacity:.5}
.grid{display:grid;grid-template-columns:repeat(2,1fr);gap:14px}
.card{background:var(--card);border-radius:18px;padding:18px;box-shadow:0 4px 12px rgba(0,0,0,.03);margin-bottom:14px}
.cf{grid-column:1/-1}
.vb{font-size:28px;font-weight:700;font-variant-numeric:tabular-nums}
.vs{font-size:14px;color:var(--sec);font-weight:500}
.vsm{font-size:12px;color:var(--sec)}
.cmini{height:90px;width:100%;margin-top:12px}
.sub{font-size:13px;color:var(--sec);font-weight:600;margin:24px 0 10px 4px;text-transform:uppercase;letter-spacing:.5px}
.li{padding:14px 0;border-bottom:.5px solid var(--brd);display:flex;justify-content:space-between;align-items:center}
.li:last-child{border-bottom:none}
.li.ck{cursor:pointer}
.li.ck:active{background:rgba(128,128,128,.1);border-radius:8px}
.ll{display:flex;align-items:center;gap:12px;flex:1;min-width:0}
.licon{width:36px;height:36px;border-radius:10px;display:flex;align-items:center;justify-content:center;font-size:18px;flex-shrink:0}
.ln{font-weight:600;font-size:15px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.lm{font-size:12px;color:var(--sec);margin-top:2px}
.lr{display:flex;align-items:center;gap:8px;flex-shrink:0}
.badge{font-size:11px;padding:3px 8px;border-radius:8px;font-weight:600}
.b-ok{background:rgba(52,199,89,.15);color:var(--ok)}
.b-err{background:rgba(255,59,48,.15);color:var(--danger)}
.b-warn{background:rgba(255,149,0,.15);color:var(--warn)}
.b-info{background:rgba(0,122,255,.1);color:var(--accent)}
.chev{color:var(--sec);font-size:14px}
.btn{padding:14px;border-radius:14px;border:none;font-weight:600;font-size:16px;cursor:pointer;width:100%;transition:.2s}
.btn:active{transform:scale(.98)}
.btn-p{background:var(--accent);color:#fff}
.btn-d{background:var(--danger);color:#fff}
.btn-s{background:rgba(128,128,128,.12);color:var(--text)}
.btn-sm{padding:8px 16px;font-size:13px;width:auto;border-radius:10px}
.fab{position:fixed;bottom:calc(100px + env(safe-area-inset-bottom));right:max(20px,calc((100vw - 600px)/2));width:56px;height:56px;border-radius:28px;background:var(--accent);color:#fff;border:none;font-size:28px;cursor:pointer;box-shadow:0 6px 20px rgba(0,122,255,.4);z-index:100;display:none;align-items:center;justify-content:center}
.fab:active{transform:scale(.92)}
.fab.show{display:flex}
.dock{position:fixed;bottom:calc(20px + env(safe-area-inset-bottom));left:50%;transform:translateX(-50%);width:92%;max-width:420px;height:60px;background:var(--dock);backdrop-filter:blur(25px);-webkit-backdrop-filter:blur(25px);border-radius:30px;display:flex;justify-content:space-evenly;align-items:center;box-shadow:0 10px 30px rgba(0,0,0,.12);z-index:1000}
.dbtn{border:none;background:none;width:40px;height:40px;border-radius:12px;display:flex;flex-direction:column;align-items:center;justify-content:center;gap:1px;cursor:pointer;color:var(--sec);transition:.2s}
.dbtn.act{color:var(--accent);background:rgba(0,122,255,.1)}
.dbtn svg{width:20px;height:20px;fill:currentColor}
.dlbl{font-size:8px;font-weight:600}
.pg{display:none;opacity:0;transition:opacity .25s ease}
.pg.act{display:block;opacity:1}
.mov{position:fixed;inset:0;background:rgba(0,0,0,.4);backdrop-filter:blur(8px);display:none;align-items:center;justify-content:center;z-index:2000;opacity:0;transition:opacity .2s}
.mov.act{display:flex;opacity:1}
.mcard{background:var(--card);width:88%;max-width:450px;border-radius:20px;padding:24px;max-height:80vh;overflow-y:auto;transform:scale(.9);transition:transform .2s}
.mov.act .mcard{transform:scale(1)}
.mtitle{font-size:20px;font-weight:700;margin-bottom:16px}
.minput{background:rgba(128,128,128,.08);border:none;padding:12px 16px;border-radius:12px;width:100%;color:var(--text);font-size:15px;margin-bottom:12px;outline:none}
.minput:focus{box-shadow:0 0 0 2px var(--accent)}
.mlabel{font-size:12px;color:var(--sec);font-weight:600;margin-bottom:6px;display:block}
.msel{background:rgba(128,128,128,.08);border:none;padding:12px 16px;border-radius:12px;width:100%;color:var(--text);font-size:15px;margin-bottom:12px;outline:none;appearance:none}
.macts{display:flex;gap:10px;margin-top:16px}
.macts .btn{flex:1}
.pbar{height:8px;background:rgba(128,128,128,.15);border-radius:4px;overflow:hidden;margin-top:6px}
.pfill{height:100%;border-radius:4px;transition:width .5s}
.empty{text-align:center;padding:40px 20px;color:var(--sec)}
.eicon{font-size:48px;margin-bottom:12px}
.abanner{background:rgba(255,59,48,.12);border:1px solid var(--danger);color:var(--danger);border-radius:14px;padding:14px 18px;margin-bottom:14px;display:none;align-items:center;gap:10px;font-size:13px;font-weight:600}
.abanner.show{display:flex}
.sbar{background:rgba(128,128,128,.12);border:none;padding:10px 16px;border-radius:12px;width:100%;color:var(--text);font-size:15px;outline:none;margin-bottom:14px}
.sbar::placeholder{color:var(--sec)}
.seg{display:flex;background:rgba(128,128,128,.1);padding:2px;border-radius:10px;margin-bottom:14px}
.segi{flex:1;border:none;background:none;color:var(--sec);padding:8px 4px;border-radius:8px;font-size:12px;font-weight:600;cursor:pointer;text-align:center;transition:.2s}
.segi.act{background:var(--card);color:var(--text);box-shadow:0 2px 6px rgba(0,0,0,.08)}
.toast{position:fixed;top:calc(20px + env(safe-area-inset-top));left:50%;transform:translateX(-50%) translateY(-100px);background:var(--card);color:var(--text);padding:12px 24px;border-radius:14px;font-weight:600;font-size:14px;box-shadow:0 8px 30px rgba(0,0,0,.15);z-index:3000;transition:transform .3s;pointer-events:none}
.toast.show{transform:translateX(-50%) translateY(0)}
.logbox{background:var(--bg);padding:14px;border-radius:12px;font-family:'SF Mono',monospace;font-size:11px;white-space:pre-wrap;max-height:300px;overflow-y:auto;line-height:1.6}
.srow{display:flex;justify-content:space-between;padding:8px 0}
.slbl{color:var(--sec);font-size:13px}
.sval{font-weight:600;font-size:13px}
.loader{text-align:center;padding:40px;color:var(--sec)}
.spinner{width:30px;height:30px;border:3px solid var(--brd);border-top-color:var(--accent);border-radius:50%;animation:spin 1s linear infinite;margin:0 auto 12px}
@keyframes spin{to{transform:rotate(360deg)}}
.api-item{background:var(--bg);border-radius:10px;padding:12px;margin-bottom:8px;font-size:12px}
.api-method{display:inline-block;padding:2px 8px;border-radius:4px;font-weight:700;font-size:10px;margin-right:8px}
.api-get{background:rgba(0,122,255,.15);color:var(--accent)}
.api-post{background:rgba(52,199,89,.15);color:var(--ok)}
.api-put{background:rgba(255,149,0,.15);color:var(--warn)}
.api-delete{background:rgba(255,59,48,.15);color:var(--danger)}
.api-path{font-family:'SF Mono',monospace;font-weight:600}
.api-desc{color:var(--sec);margin-top:4px;font-size:11px}
.chart-lg{height:200px;width:100%;margin-top:12px}
</style>
</head>
<body>
<div class="toast" id="toast"></div>
<div class="mov" id="loginModal"><div class="mcard"><div class="mtitle">登录 iSCSI Panel</div><label class="mlabel">用户名</label><input class="minput" id="loginUser" placeholder="admin" value="admin"><label class="mlabel">密码</label><input class="minput" id="loginPass" type="password" placeholder="password" value="admin123"><div class="macts"><button class="btn btn-p" onclick="doLogin()">登录</button></div></div></div>
<div class="mov" id="deleteModal"><div class="mcard"><div class="mtitle">确认删除</div><p style="color:var(--sec);font-size:14px;margin-bottom:16px" id="deleteMsg">确定要删除吗？此操作不可撤销。</p><div class="macts"><button class="btn btn-s" onclick="closeModal('deleteModal')">取消</button><button class="btn btn-d" id="deleteConfirmBtn">删除</button></div></div></div>
<div class="mov" id="formModal"><div class="mcard"><div class="mtitle" id="formTitle">创建</div><div id="formBody"></div><div class="macts"><button class="btn btn-s" onclick="closeModal('formModal')">取消</button><button class="btn btn-p" id="formSubmitBtn">保存</button></div></div></div>
<div class="ctn">
<div class="hdr"><div><div class="hdr-sub" id="dateNow"></div><div class="ttl" id="pgTitle">仪表盘</div></div><div class="hdr-r"><div class="sdot">● 运行中</div><div class="utime" id="uptime">--</div></div></div>
<div class="abanner" id="alertBanner"><span>⚠️</span><span id="alertMsg" style="flex:1"></span></div>

<div id="p-dash" class="pg act">
<div class="grid">
<div class="card"><div class="vs" style="color:var(--accent)">Target 数量</div><div class="vb" id="dTargets">-</div><div class="vsm">活跃 <span id="dTargetsAct">-</span></div></div>
<div class="card"><div class="vs" style="color:var(--purple)">LUN 数量</div><div class="vb" id="dLuns">-</div><div class="vsm">总容量 <span id="dLunsSize">-</span></div></div>
<div class="card"><div class="vs" style="color:var(--ok)">连接数</div><div class="vb" id="dConn">-</div><div class="vsm">存储 <span id="dStorage">-</span></div></div>
<div class="card"><div class="vs" style="color:var(--warn)">告警</div><div class="vb" id="dAlerts">-</div><div class="vsm">运行 <span id="dUptime">-</span></div></div>
</div>
<div class="card"><div style="display:flex;justify-content:space-between;margin-bottom:5px"><span class="vs">系统资源</span><span style="font-size:12px;font-weight:600" id="dCpuTxt">-</span></div><div class="cmini"><canvas id="cpuChart"></canvas></div></div>
<div class="card"><div style="display:flex;justify-content:space-between;margin-bottom:5px"><span class="vs">吞吐量</span><span style="font-size:12px;font-weight:600"><span style="color:var(--ok)">↓ <span id="