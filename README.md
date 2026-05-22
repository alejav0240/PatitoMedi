# PatitoMedi


Perfecto, entendido completamente. Déjame diseñar la arquitectura completa. Los componentes son claros, voy a plasmarlo como un diagrama interactivo — haz clic en cualquier caja para profundizar en ese componente.


<style>
*{box-sizing:border-box;margin:0;padding:0}
.wrap{padding:20px 16px 24px;position:relative;font-family:var(--font-sans)}
#sv{position:absolute;inset:0;width:100%;pointer-events:none;overflow:visible}
.row{display:flex;justify-content:center;gap:7px}
.row.s>*{flex:1;min-width:0}
.sp{height:30px}.sps{height:11px}
.bx{border-radius:8px;padding:7px 8px;border:0.5px solid;text-align:center;cursor:pointer;transition:opacity .15s}
.bx:hover{opacity:.8}
.bt{font-size:11px;font-weight:500}
.bs{font-size:10px;opacity:.85;margin-top:2px}
.tag{display:inline-block;font-size:9px;padding:1px 4px;border-radius:3px;margin-top:3px;font-weight:500}
.lbl{font-size:10px;font-weight:500;color:var(--color-text-tertiary);text-transform:uppercase;letter-spacing:.06em;margin-bottom:5px;padding-left:2px}
.infra-row{display:flex;gap:8px;justify-content:center;align-items:stretch}
.infra-row .side{display:flex;flex-direction:column;gap:7px}
.gray{background:var(--color-background-secondary);border-color:var(--color-border-secondary);color:var(--color-text-primary)}
.blue{background:#E6F1FB;border-color:#185FA5;color:#0C447C}
.pur{background:#EEEDFE;border-color:#534AB7;color:#3C3489}
.tea{background:#E1F5EE;border-color:#0F6E56;color:#085041}
.cor{background:#FAECE7;border-color:#993C1D;color:#712B13}
.amb{background:#FAEEDA;border-color:#854F0B;color:#633806}
.pnk{background:#FBEAF0;border-color:#993556;color:#72243E}
.grn{background:#EAF3DE;border-color:#3B6D11;color:#27500A}
.kfk{background:#3C3489;border-color:#7F77DD;color:#EEEDFE}
.cya{background:#E1F5EE;border-color:#1D9E75;color:#0F6E56}
.red{background:#FCEBEB;border-color:#A32D2D;color:#791F1F}
.trest{background:#E6F1FB;color:#185FA5}
.tgql{background:#FBEAF0;color:#993556}
.tws{background:#EAF3DE;color:#3B6D11}
@media(prefers-color-scheme:dark){
.blue{background:#0C447C;border-color:#85B7EB;color:#B5D4F4}
.pur{background:#3C3489;border-color:#AFA9EC;color:#CECBF6}
.tea{background:#085041;border-color:#5DCAA5;color:#9FE1CB}
.cor{background:#712B13;border-color:#F0997B;color:#F5C4B3}
.amb{background:#633806;border-color:#EF9F27;color:#FAC775}
.pnk{background:#72243E;border-color:#ED93B1;color:#F4C0D1}
.grn{background:#27500A;border-color:#97C459;color:#C0DD97}
.kfk{background:#26215C;border-color:#7F77DD;color:#CECBF6}
.cya{background:#085041;border-color:#5DCAA5;color:#9FE1CB}
.red{background:#791F1F;border-color:#F09595;color:#F7C1C1}
.trest{background:#0C447C;color:#B5D4F4}
.tgql{background:#72243E;color:#F4C0D1}
.tws{background:#085041;color:#9FE1CB}
}
</style>

<div class="wrap" id="arch">
<svg id="sv"></svg>

<div class="lbl">Clientes</div>
<div class="row">
  <div class="bx gray" id="bweb" onclick="sendPrompt('¿Cómo se conectan la web y mobile app al gateway? ¿Cómo manejan la señalización WebRTC desde el cliente?')">
    <div class="bt">Web App</div><div class="bs">React / Vue</div>
  </div>
  <div class="bx gray" id="bmob" onclick="sendPrompt('¿Cómo se conectan la web y mobile app al gateway? ¿Cómo manejan la señalización WebRTC desde el cliente?')">
    <div class="bt">Mobile App</div><div class="bs">iOS · Android</div>
  </div>
</div>

<div class="sp"></div>

<div class="lbl">Entrada</div>
<div class="row">
  <div class="bx blue" id="bnginx" onclick="sendPrompt('Dame la config Nginx para telemedicina incluyendo el upgrade de WebSocket para el servicio de video llamada.')">
    <div class="bt">Nginx</div><div class="bs">Reverse proxy · SSL · Load balancing · WS Upgrade</div>
  </div>
</div>
<div class="sps"></div>
<div class="row">
  <div class="bx pur" id="bkong" onclick="sendPrompt('Config Kong para telemedicina: JWT, rate limiting, CORS, GraphQL proxy, y WebSocket passthrough para el video service.')">
    <div class="bt">Kong API Gateway</div>
    <div class="bs">JWT · Rate limiting · REST + GraphQL + WebSocket routing</div>
  </div>
</div>

<div class="sp"></div>

<div class="lbl">Microservicios</div>
<div class="row s">
  <div class="bx tea" id="buser" onclick="sendPrompt('Diseña el User Service en Golang: endpoints REST, JWT, esquema PostgreSQL para pacientes y médicos, y eventos Kafka que publica.')">
    <div class="bt">User Service</div>
    <div class="bs">Golang</div>
    <div class="bs">Auth · Profiles</div>
    <div><span class="tag trest">REST</span></div>
  </div>
  <div class="bx cor" id="bappt" onclick="sendPrompt('Diseña el Appointments Service en Spring Boot: scheduling, slots, esquema PostgreSQL, y eventos Kafka.')">
    <div class="bt">Appointments</div>
    <div class="bs">Spring Boot</div>
    <div class="bs">Scheduling</div>
    <div><span class="tag trest">REST</span></div>
  </div>
  <div class="bx amb" id="bpay" onclick="sendPrompt('Diseña el Payments Service en Django: billing, invoices, PostgreSQL, y eventos Kafka.')">
    <div class="bt">Payments</div>
    <div class="bs">Django</div>
    <div class="bs">Billing</div>
    <div><span class="tag trest">REST</span></div>
  </div>
  <div class="bx pnk" id="bmed" onclick="sendPrompt('Diseña el Medical History Service en Express.js: GraphQL schema completo, resolvers, colecciones MongoDB.')">
    <div class="bt">Medical Hist.</div>
    <div class="bs">Express.js</div>
    <div class="bs">EHR · Docs</div>
    <div><span class="tag tgql">GraphQL</span></div>
  </div>
  <div class="bx cya" id="bvideo" onclick="sendPrompt('Diseña en detalle el Video Call Service en Golang con WebRTC: arquitectura de señalización WebSocket (offer/answer/ICE), integración con coturn STUN/TURN, manejo de salas de videoconferencia, Redis para sesiones, y eventos Kafka call-started y call-ended.')">
    <div class="bt">Video Call</div>
    <div class="bs">Golang</div>
    <div class="bs">Signaling · Rooms</div>
    <div><span class="tag tws">WebSocket</span></div>
  </div>
</div>

<div class="sps"></div>

<!-- STUN/TURN row, aligned to video call -->
<div class="row s" style="justify-content:flex-end">
  <div style="flex:0 0 calc(20% - 3px)">
    <div class="bx gray" id="bstun" onclick="sendPrompt('Cómo configurar coturn como servidor STUN/TURN para WebRTC en la plataforma de telemedicina? Include autenticación con HMAC, puertos, y configuración de relay.')">
      <div class="bt" style="font-size:10px">STUN/TURN</div>
      <div class="bs">coturn</div>
    </div>
  </div>
</div>

<div class="sp"></div>

<div class="lbl">Mensajería asíncrona</div>
<div class="row">
  <div class="bx kfk" id="bkafka" onclick="sendPrompt('Define todos los Kafka topics para telemedicina incluyendo los nuevos: call-started, call-ended, call-recording-ready. Para cada topic: productor, consumidores, payload, particiones y retención.')">
    <div class="bt">Apache Kafka — Event Bus</div>
    <div class="bs" style="opacity:.75">appointment-created · payment-confirmed · user-registered · record-updated · call-started · call-ended</div>
  </div>
</div>

<div class="sps"></div>

<div class="lbl">Capa de datos</div>
<div class="row s">
  <div class="bx blue" id="bdbu" onclick="sendPrompt('Esquema PostgreSQL para usuarios: tablas patients, doctors, sessions, roles.')">
    <div class="bt">PostgreSQL 1</div><div class="bs">Users DB</div>
  </div>
  <div class="bx blue" id="bdba" onclick="sendPrompt('Esquema PostgreSQL para citas: appointments, slots, schedules.')">
    <div class="bt">PostgreSQL 2</div><div class="bs">Appointments DB</div>
  </div>
  <div class="bx blue" id="bdbp" onclick="sendPrompt('Esquema PostgreSQL para pagos: transactions, invoices, audit_log.')">
    <div class="bt">PostgreSQL 3</div><div class="bs">Payments DB</div>
  </div>
  <div class="bx grn" id="bdbm" onclick="sendPrompt('Colecciones MongoDB para historial médico: clinical_records, prescriptions, attachments.')">
    <div class="bt">MongoDB</div><div class="bs">Medical Records</div>
  </div>
  <div class="bx red" id="bdbredis" onclick="sendPrompt('¿Cómo usar Redis en el Video Call Service? Define la estructura de datos para: sesiones WebRTC activas, rooms, participantes, y TTL de cada clave.')">
    <div class="bt">Redis</div><div class="bs">Call Sessions</div>
  </div>
</div>

<div class="sp"></div>

<div class="lbl">Monitoreo</div>
<div class="row" style="gap:16px">
  <div class="bx cor" id="bprom" style="min-width:120px" onclick="sendPrompt('¿Qué métricas Prometheus exponer en cada servicio? Incluye métricas específicas para el Video Call Service: duración de llamadas, peers conectados, ICE failures.')">
    <div class="bt">Prometheus</div><div class="bs">Metrics · Scraping</div>
  </div>
  <div class="bx amb" id="bgraf" style="min-width:120px" onclick="sendPrompt('Dashboards Grafana para telemedicina: paneles de salud de servicios, latencia p99, video call concurrency, y alertas de negocio.')">
    <div class="bt">Grafana</div><div class="bs">Dashboards · Alerts</div>
  </div>
</div>

<div style="text-align:center;font-size:10px;color:var(--color-text-tertiary);margin-top:16px">
  Database-per-Service pattern. Kafka gestiona eventos asíncronos. WebRTC usa coturn para atravesar NAT/firewalls. Haz clic para profundizar.
</div>
</div>

<script>
(function(){
const arch=document.getElementById('arch');
const sv=document.getElementById('sv');
function rc(id){
  const el=document.getElementById(id);if(!el)return null;
  const e=el.getBoundingClientRect(),a=arch.getBoundingClientRect();
  return{t:e.top-a.top,b:e.bottom-a.top,l:e.left-a.left,r:e.right-a.left,
    cx:(e.left+e.right)/2-a.left,cy:(e.top+e.bottom)/2-a.top,w:e.width,h:e.height};
}
function draw(){
  sv.innerHTML='';
  const defs=document.createElementNS('http://www.w3.org/2000/svg','defs');
  defs.innerHTML=`<marker id="at" viewBox="0 0 10 10" refX="8" refY="5" markerWidth="5" markerHeight="5" orient="auto-start-reverse"><path d="M2 1L8 5L2 9" fill="none" stroke="context-stroke" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></marker>`;
  sv.appendChild(defs);
  sv.setAttribute('height',(arch.scrollHeight+10)+'px');

  const web=rc('bweb'),mob=rc('bmob'),nginx=rc('bnginx'),kong=rc('bkong');
  const user=rc('buser'),appt=rc('bappt'),pay=rc('bpay'),med=rc('bmed'),video=rc('bvideo');
  const stun=rc('bstun'),kafka=rc('bkafka');
  const dbu=rc('bdbu'),dba=rc('bdba'),dbp=rc('bdbp'),dbm=rc('bdbm'),redis=rc('bdbredis');
  const prom=rc('bprom'),graf=rc('bgraf');
  if(!web||!nginx)return;

  function pa(x1,y1,x2,y2,o){
    o=o||{};
    const el=document.createElementNS('http://www.w3.org/2000/svg','path');
    const my=(y1+y2)/2;
    const d=(Math.abs(x2-x1)<4)?`M${x1} ${y1}L${x2} ${y2}`:`M${x1} ${y1}C${x1} ${my},${x2} ${my},${x2} ${y2}`;
    el.setAttribute('d',d);el.setAttribute('fill','none');
    el.setAttribute('stroke',o.c||'var(--color-border-primary)');
    el.setAttribute('stroke-width',o.w||'1');
    el.setAttribute('opacity',o.op||'0.5');
    if(o.dash)el.setAttribute('stroke-dasharray',o.dash);
    el.setAttribute('marker-end','url(#at)');
    sv.appendChild(el);
    if(o.lbl){
      const t=document.createElementNS('http://www.w3.org/2000/svg','text');
      t.setAttribute('x',(x1+x2)/2+5);t.setAttribute('y',my-2);
      t.setAttribute('font-size','9');t.setAttribute('font-family','var(--font-sans)');
      t.setAttribute('fill',o.lc||'var(--color-text-tertiary)');
      t.setAttribute('dominant-baseline','central');t.textContent=o.lbl;
      sv.appendChild(t);
    }
  }

  // clients → nginx
  pa(web.cx,web.b,nginx.cx,nginx.t);
  pa(mob.cx,mob.b,nginx.cx,nginx.t);
  // nginx → kong
  pa(nginx.cx,nginx.b,kong.cx,kong.t);
  // kong → services
  [user,appt,pay].forEach(s=>pa(kong.cx,kong.b,s.cx,s.t,{op:'0.4'}));
  pa(kong.cx,kong.b,med.cx,med.t,{op:'0.4',lbl:'GraphQL',lc:'#993556'});
  pa(kong.cx,kong.b,video.cx,video.t,{op:'0.4',lbl:'WS',lc:'#0F6E56'});

  // video → stun/turn (short downward)
  if(stun) pa(video.cx,video.b,stun.cx,stun.t,{c:'#1D9E75',dash:'3 2',op:'0.6',w:'0.9'});

  // services → kafka (dashed)
  [user,appt,pay,med,video].forEach(s=>pa(s.cx,s.b > (stun?stun.b:s.b) ? s.b : s.b, s.cx, kafka.t, {c:'#7F77DD',dash:'4 3',op:'0.55'}));
  // stun connects to kafka via video (not directly — skip)

  // kafka → dbs (aligned columns)
  [[user,dbu],[appt,dba],[pay,dbp],[med,dbm],[video,redis]].forEach(([s,db])=>{
    if(db)pa(s.cx,kafka.b,db.cx,db.t,{op:'0.3',w:'0.8',dash:'3 2'});
  });

  // prometheus → grafana
  if(prom&&graf)pa(prom.r,prom.cy,graf.l,graf.cy,{op:'0.6'});

  // prometheus scrapes (dashed up)
  if(prom&&kafka){
    const el=document.createElementNS('http://www.w3.org/2000/svg','path');
    const tY=kafka.b+5;
    el.setAttribute('d',`M${prom.cx} ${prom.t} L${prom.cx} ${tY} L${kafka.r-6} ${tY}`);
    el.setAttribute('fill','none');el.setAttribute('stroke','var(--color-border-secondary)');
    el.setAttribute('stroke-width','0.8');el.setAttribute('stroke-dasharray','3 3');
    el.setAttribute('opacity','0.35');sv.appendChild(el);
    const t=document.createElementNS('http://www.w3.org/2000/svg','text');
    t.setAttribute('x',prom.cx+7);t.setAttribute('y',(prom.t*0.6+tY*0.4));
    t.setAttribute('font-size','9');t.setAttribute('font-family','var(--font-sans)');
    t.setAttribute('fill','var(--color-text-tertiary)');
    t.setAttribute('dominant-baseline','central');t.textContent='scrapes /metrics';
    sv.appendChild(t);
  }
}
setTimeout(draw,150);
window.addEventListener('resize',()=>setTimeout(draw,50));
})();
</script>
