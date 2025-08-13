package frame

const (
	// Html401 StatusUnauthorized
	Html401 = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>401 - 未授权</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Microsoft YaHei', Arial, sans-serif;
            background-color: #f0f4f0;
            color: #333;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            overflow: hidden;
            position: relative;
        }
        
        .background {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            z-index: -1;
            overflow: hidden;
        }
        
        .circle {
            position: absolute;
            border-radius: 50%;
            background: linear-gradient(45deg, rgba(79,70,229,0.1), rgba(129,140,248,0.2));
            animation: float var(--duration) infinite ease-in-out;
            opacity: 0.6;
        }
        
        @keyframes float {
            0%, 100% { transform: translateY(0) translateX(0) rotate(0deg); }
            25% { transform: translateY(-20px) translateX(10px) rotate(5deg); }
            50% { transform: translateY(0) translateX(20px) rotate(10deg); }
            75% { transform: translateY(20px) translateX(10px) rotate(5deg); }
        }
        
        .container {
            text-align: center;
            background: linear-gradient(135deg, rgba(255, 255, 255, 0.9), rgba(240, 245, 255, 0.85));
            backdrop-filter: blur(10px);
            border-radius: 30px;
            padding: 50px;
            box-shadow: 0 20px 50px rgba(0, 0, 0, 0.1), inset 0 0 30px rgba(79, 70, 229, 0.1);
            max-width: 90%;
            width: 650px;
            position: relative;
            z-index: 1;
            animation: pulse 6s infinite ease-in-out;
            border: 1px solid rgba(255, 255, 255, 0.5);
            overflow: hidden;
        }
        
        @keyframes pulse {
            0%, 100% { transform: scale(1); }
            50% { transform: scale(1.02); }
        }
        
        .error-code {
            font-size: 180px;
            font-weight: 900;
            margin-bottom: 0;
            background: linear-gradient(45deg, #4f46e5, #818cf8, #4f46e5);
            background-size: 200% 200%;
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            line-height: 1;
            text-shadow: 0 5px 15px rgba(79, 70, 229, 0.2);
            animation: gradient 8s ease infinite;
            letter-spacing: -5px;
            position: relative;
        }
        
        @keyframes gradient {
            0% { background-position: 0% 50%; }
            50% { background-position: 100% 50%; }
            100% { background-position: 0% 50%; }
        }
        
        .error-text {
            font-size: 32px;
            font-weight: 700;
            margin: 10px 0 30px;
            color: #4338ca;
            position: relative;
            display: inline-block;
            padding: 0 15px;
        }
        
        .error-text::before,
        .error-text::after {
            content: '';
            position: absolute;
            height: 2px;
            width: 30px;
            background: linear-gradient(45deg, #4f46e5, #818cf8);
            top: 50%;
            transform: translateY(-50%);
        }
        
        .error-text::before {
            left: -40px;
        }
        
        .error-text::after {
            right: -40px;
        }
        
        .error-description {
            font-size: 20px;
            color: #4a5568;
            margin-bottom: 20px;
            line-height: 1.8;
            background: linear-gradient(45deg, #4a5568, #718096);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            padding: 10px 20px;
            position: relative;
            z-index: 2;
            text-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
        }
        
        .key {
            width: 150px;
            height: 150px;
            margin: 0 auto 30px;
            position: relative;
            filter: drop-shadow(0 10px 20px rgba(79, 70, 229, 0.3));
        }
        
        .key svg {
            width: 100%;
            height: 100%;
            animation: float-key 3s infinite ease-in-out;
            transform-origin: center center;
        }
        
        @keyframes float-key {
            0%, 100% { transform: translateY(0) rotate(0deg); }
            50% { transform: translateY(-10px) rotate(5deg); }
        }
        
        .badge {
            position: absolute;
            width: 40px;
            height: 40px;
            background: linear-gradient(45deg, #4f46e5, #818cf8);
            border-radius: 50%;
            display: flex;
            justify-content: center;
            align-items: center;
            animation: pulse-badge 2s infinite ease-in-out;
            box-shadow: 0 5px 15px rgba(79, 70, 229, 0.3);
        }
        
        .badge::before {
            content: '?';
            color: white;
            font-weight: bold;
            font-size: 24px;
        }
        
        @keyframes pulse-badge {
            0%, 100% { transform: scale(1); }
            50% { transform: scale(1.1); }
        }
        
        .badge-1 {
            top: 10%;
            left: 10%;
            animation-delay: 0s;
        }
        
        .badge-2 {
            bottom: 15%;
            right: 10%;
            animation-delay: 0.5s;
        }
        
        .badge-3 {
            top: 20%;
            right: 15%;
            animation-delay: 1s;
        }
        
        .particles {
            position: absolute;
            width: 100%;
            height: 100%;
            top: 0;
            left: 0;
            pointer-events: none;
        }
        
        .particle {
            position: absolute;
            width: 8px;
            height: 8px;
            background-color: rgba(129, 140, 248, 0.5);
            border-radius: 50%;
            animation: float-particle 15s infinite linear;
        }
        
        @keyframes float-particle {
            0% {
                transform: translateY(0) rotate(0deg);
                opacity: 0;
            }
            10% {
                opacity: 1;
            }
            90% {
                opacity: 1;
            }
            100% {
                transform: translateY(-100px) rotate(360deg);
                opacity: 0;
            }
        }
        
        .wave {
            position: absolute;
            width: 100%;
            height: 20px;
            background: linear-gradient(90deg, transparent, rgba(79, 70, 229, 0.3), transparent);
            opacity: 0.5;
            animation: wave 8s linear infinite;
        }
        
        .wave-top {
            top: 15%;
            animation-delay: 0s;
        }
        
        .wave-middle {
            top: 50%;
            animation-delay: 2s;
        }
        
        .wave-bottom {
            bottom: 15%;
            animation-delay: 4s;
        }
        
        @keyframes wave {
            0% { transform: translateX(-100%) scaleY(1); }
            50% { transform: translateX(0%) scaleY(0.5); }
            100% { transform: translateX(100%) scaleY(1); }
        }
        
        @media (max-width: 768px) {
            .error-code {
                font-size: 120px;
            }
            
            .error-text {
                font-size: 24px;
            }
            
            .container {
                padding: 30px;
            }
            
            .key {
                width: 100px;
                height: 100px;
            }
            
            .error-text::before,
            .error-text::after {
                width: 20px;
            }
            
            .error-text::before {
                left: -25px;
            }
            
            .error-text::after {
                right: -25px;
            }
            
            .badge {
                width: 30px;
                height: 30px;
            }
            
            .badge::before {
                font-size: 18px;
            }
        }
    </style>
</head>
<body>
    <div class="background" id="background"></div>
    
    <div class="badge badge-1"></div>
    <div class="badge badge-2"></div>
    <div class="badge badge-3"></div>
    
    <div class="wave wave-top"></div>
    <div class="wave wave-middle"></div>
    <div class="wave wave-bottom"></div>
    
    <div class="container">
        <div class="key">
            <svg viewBox="0 0 512 512" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M336 192h-16v-64C320 57.31 262.69 0 192 0S64 57.31 64 128v64H48c-26.51 0-48 21.49-48 48v192c0 26.51 21.49 48 48 48h288c26.51 0 48-21.49 48-48V240c0-26.51-21.49-48-48-48zM192 64c35.29 0 64 28.71 64 64v64H128v-64c0-35.29 28.71-64 64-64zm112 352H96c-8.84 0-16-7.16-16-16s7.16-16 16-16h208c8.84 0 16 7.16 16 16s-7.16 16-16 16zm32-80c0 8.84-7.16 16-16 16H96c-8.84 0-16-7.16-16-16v-48c0-8.84 7.16-16 16-16h224c8.84 0 16 7.16 16 16v48z" fill="#4F46E5"/>
            </svg>
        </div>
        <h1 class="error-code">401</h1>
        <h2 class="error-text">未授权</h2>
        <p class="error-description">抱歉！您需要登录或提供有效凭证才能访问此页面。<br>请尝试登录或联系管理员获取访问权限。</p>
        <div class="particles" id="particles"></div>
    </div>

    <script>
        // 创建圆形背景
        function createCircles() {
            const background = document.getElementById('background');
            const count = 15;
            
            for (let i = 0; i < count; i++) {
                const circle = document.createElement('div');
                circle.className = 'circle';
                
                // 随机位置
                const x = Math.random() * 100;
                const y = Math.random() * 100;
                circle.style.left = x + '%';
                circle.style.top = y + '%';
                
                // 随机大小
                const size = 50 + Math.random() * 100;
                circle.style.width = size + 'px';
                circle.style.height = size + 'px';
                
                // 随机动画持续时间
                const duration = 10 + Math.random() * 20;
                circle.style.setProperty('--duration', duration + 's');
                
                // 随机延迟
                const delay = Math.random() * 10;
                circle.style.animationDelay = delay + 's';
                
                background.appendChild(circle);
            }
        }
        
        // 创建粒子
        function createParticles() {
            const container = document.getElementById('particles');
            const count = 20;
            
            for (let i = 0; i < count; i++) {
                const particle = document.createElement('div');
                particle.className = 'particle';
                
                // 随机位置
                const x = Math.random() * 100;
                const y = Math.random() * 100;
                particle.style.left = x + '%';
                particle.style.top = y + '%';
                
                // 随机大小
                const size = Math.random() * 5 + 3;
                particle.style.width = size + 'px';
                particle.style.height = size + 'px';
                
                // 随机动画延迟
                const delay = Math.random() * 10;
                particle.style.animationDelay = delay + 's';
                
                container.appendChild(particle);
            }
        }
        
        // 页面加载时创建元素
        window.addEventListener('load', function() {
            createCircles();
            createParticles();
        });
    </script>
</body>
</html>`

	// Html403 StatusForbidden
	Html403 = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>403 - 禁止访问</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Microsoft YaHei', Arial, sans-serif;
            background-color: #f8f0f4;
            color: #333;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            overflow: hidden;
            position: relative;
        }
        
        .background {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            z-index: -1;
            overflow: hidden;
        }
        
        .hexagon {
            position: absolute;
            clip-path: polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%);
            background: linear-gradient(45deg, rgba(220,38,38,0.1), rgba(248,113,113,0.2));
            animation: float var(--duration) infinite ease-in-out;
            opacity: 0.6;
        }
        
        @keyframes float {
            0%, 100% { transform: translateY(0) translateX(0) rotate(0deg); }
            25% { transform: translateY(-20px) translateX(10px) rotate(5deg); }
            50% { transform: translateY(0) translateX(20px) rotate(10deg); }
            75% { transform: translateY(20px) translateX(10px) rotate(5deg); }
        }
        
        .container {
            text-align: center;
            background: linear-gradient(135deg, rgba(255, 255, 255, 0.9), rgba(255, 240, 245, 0.85));
            backdrop-filter: blur(10px);
            border-radius: 30px;
            padding: 50px;
            box-shadow: 0 20px 50px rgba(0, 0, 0, 0.1), inset 0 0 30px rgba(220, 38, 38, 0.1);
            max-width: 90%;
            width: 650px;
            position: relative;
            z-index: 1;
            animation: pulse 6s infinite ease-in-out;
            border: 1px solid rgba(255, 255, 255, 0.5);
            overflow: hidden;
        }
        
        @keyframes pulse {
            0%, 100% { transform: scale(1); }
            50% { transform: scale(1.02); }
        }
        
        .error-code {
            font-size: 180px;
            font-weight: 900;
            margin-bottom: 0;
            background: linear-gradient(45deg, #dc2626, #f87171, #dc2626);
            background-size: 200% 200%;
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            line-height: 1;
            text-shadow: 0 5px 15px rgba(220, 38, 38, 0.2);
            animation: gradient 8s ease infinite;
            letter-spacing: -5px;
            position: relative;
        }
        
        @keyframes gradient {
            0% { background-position: 0% 50%; }
            50% { background-position: 100% 50%; }
            100% { background-position: 0% 50%; }
        }
        
        .error-text {
            font-size: 32px;
            font-weight: 700;
            margin: 10px 0 30px;
            color: #991b1b;
            position: relative;
            display: inline-block;
            padding: 0 15px;
        }
        
        .error-text::before,
        .error-text::after {
            content: '';
            position: absolute;
            height: 2px;
            width: 30px;
            background: linear-gradient(45deg, #dc2626, #f87171);
            top: 50%;
            transform: translateY(-50%);
        }
        
        .error-text::before {
            left: -40px;
        }
        
        .error-text::after {
            right: -40px;
        }
        
        .error-description {
            font-size: 20px;
            color: #4a5568;
            margin-bottom: 20px;
            line-height: 1.8;
            background: linear-gradient(45deg, #4a5568, #718096);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            padding: 10px 20px;
            position: relative;
            z-index: 2;
            text-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
        }
        
        .lock {
            width: 150px;
            height: 150px;
            margin: 0 auto 30px;
            position: relative;
            filter: drop-shadow(0 10px 20px rgba(220, 38, 38, 0.3));
        }
        
        .lock svg {
            width: 100%;
            height: 100%;
            animation: shake 3s infinite ease-in-out;
            transform-origin: center bottom;
        }
        
        @keyframes shake {
            0%, 100% { transform: rotate(0deg); }
            10% { transform: rotate(-5deg); }
            20% { transform: rotate(5deg); }
            30% { transform: rotate(-5deg); }
            40% { transform: rotate(5deg); }
            50% { transform: rotate(0deg); }
        }
        
        .shield {
            position: absolute;
            width: 40px;
            height: 50px;
            background: linear-gradient(45deg, #dc2626, #f87171);
            clip-path: polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%);
            display: flex;
            justify-content: center;
            align-items: center;
            animation: float-shield 10s infinite ease-in-out;
        }
        
        .shield::before {
            content: '!';
            color: white;
            font-weight: bold;
            font-size: 24px;
        }
        
        .shield-1 {
            top: 10%;
            left: 10%;
            animation-delay: 0s;
        }
        
        .shield-2 {
            bottom: 15%;
            right: 10%;
            animation-delay: 2s;
        }
        
        .shield-3 {
            top: 20%;
            right: 15%;
            animation-delay: 4s;
        }
        
        @keyframes float-shield {
            0%, 100% { transform: translateY(0) rotate(0deg); }
            25% { transform: translateY(-15px) rotate(5deg); }
            50% { transform: translateY(0) rotate(0deg); }
            75% { transform: translateY(15px) rotate(-5deg); }
        }
        
        .particles {
            position: absolute;
            width: 100%;
            height: 100%;
            top: 0;
            left: 0;
            pointer-events: none;
        }
        
        .particle {
            position: absolute;
            width: 8px;
            height: 8px;
            background-color: rgba(248, 113, 113, 0.5);
            border-radius: 50%;
            animation: float-particle 15s infinite linear;
        }
        
        @keyframes float-particle {
            0% {
                transform: translateY(0) rotate(0deg);
                opacity: 0;
            }
            10% {
                opacity: 1;
            }
            90% {
                opacity: 1;
            }
            100% {
                transform: translateY(-100px) rotate(360deg);
                opacity: 0;
            }
        }
        
        .barrier {
            position: absolute;
            width: 100%;
            height: 10px;
            background: repeating-linear-gradient(45deg, #dc2626, #dc2626 10px, #f87171 10px, #f87171 20px);
            opacity: 0.7;
            animation: slide 20s linear infinite;
        }
        
        .barrier-top {
            top: 20%;
        }
        
        .barrier-bottom {
            bottom: 20%;
            animation-direction: reverse;
        }
        
        @keyframes slide {
            0% { transform: translateX(-100%); }
            100% { transform: translateX(100%); }
        }
        
        @media (max-width: 768px) {
            .error-code {
                font-size: 120px;
            }
            
            .error-text {
                font-size: 24px;
            }
            
            .container {
                padding: 30px;
            }
            
            .lock {
                width: 100px;
                height: 100px;
            }
            
            .error-text::before,
            .error-text::after {
                width: 20px;
            }
            
            .error-text::before {
                left: -25px;
            }
            
            .error-text::after {
                right: -25px;
            }
            
            .shield {
                width: 30px;
                height: 40px;
            }
            
            .shield::before {
                font-size: 18px;
            }
        }
    </style>
</head>
<body>
    <div class="background" id="background"></div>
    
    <div class="shield shield-1"></div>
    <div class="shield shield-2"></div>
    <div class="shield shield-3"></div>
    
    <div class="barrier barrier-top"></div>
    <div class="barrier barrier-bottom"></div>
    
    <div class="container">
        <div class="lock">
            <svg viewBox="0 0 512 512" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M336 208V128C336 57.31 278.69 0 208 0S80 57.31 80 128v80H64C28.65 208 0 236.65 0 272v192c0 35.35 28.65 64 64 64h288c35.35 0 64-28.65 64-64V272c0-35.35-28.65-64-64-64h-16zm-64 0H144V128c0-35.29 28.71-64 64-64s64 28.71 64 64v80zm-16 180c0 17.67-14.33 32-32 32s-32-14.33-32-32v-64c0-17.67 14.33-32 32-32s32 14.33 32 32v64z" fill="#DC2626"/>
            </svg>
        </div>
        <h1 class="error-code">403</h1>
        <h2 class="error-text">禁止访问</h2>
        <p class="error-description">抱歉！您没有权限访问此页面。<br>请确认您的访问权限或联系管理员。</p>
        <div class="particles" id="particles"></div>
    </div>

    <script>
        // 创建六边形背景
        function createHexagons() {
            const background = document.getElementById('background');
            const count = 15;
            
            for (let i = 0; i < count; i++) {
                const hexagon = document.createElement('div');
                hexagon.className = 'hexagon';
                
                // 随机位置
                const x = Math.random() * 100;
                const y = Math.random() * 100;
                hexagon.style.left = x + '%';
                hexagon.style.top = y + '%';
                
                // 随机大小
                const size = 50 + Math.random() * 100;
                hexagon.style.width = size + 'px';
                hexagon.style.height = size + 'px';
                
                // 随机动画持续时间
                const duration = 10 + Math.random() * 20;
                hexagon.style.setProperty('--duration', duration + 's');
                
                // 随机延迟
                const delay = Math.random() * 10;
                hexagon.style.animationDelay = delay + 's';
                
                background.appendChild(hexagon);
            }
        }
        
        // 创建粒子
        function createParticles() {
            const container = document.getElementById('particles');
            const count = 20;
            
            for (let i = 0; i < count; i++) {
                const particle = document.createElement('div');
                particle.className = 'particle';
                
                // 随机位置
                const x = Math.random() * 100;
                const y = Math.random() * 100;
                particle.style.left = x + '%';
                particle.style.top = y + '%';
                
                // 随机大小
                const size = Math.random() * 5 + 3;
                particle.style.width = size + 'px';
                particle.style.height = size + 'px';
                
                // 随机动画延迟
                const delay = Math.random() * 10;
                particle.style.animationDelay = delay + 's';
                
                container.appendChild(particle);
            }
        }
        
        // 页面加载时创建元素
        window.addEventListener('load', function() {
            createHexagons();
            createParticles();
        });
    </script>
</body>
</html>`

	// Html404 StatusNotFound
	Html404 = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>404 - 页面未找到</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Microsoft YaHei', Arial, sans-serif;
            background-color: #f8f9fa;
            color: #333;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            overflow: hidden;
            position: relative;
        }
        
        .stars {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            z-index: -1;
            overflow: hidden;
        }
        
        .star {
            position: absolute;
            background-color: #fff;
            border-radius: 50%;
            animation: twinkle var(--duration) infinite ease-in-out;
            opacity: 0;
        }
        
        @keyframes twinkle {
            0%, 100% { opacity: 0; }
            50% { opacity: var(--opacity); }
        }
        
        .container {
            text-align: center;
            background: linear-gradient(135deg, rgba(255, 255, 255, 0.9), rgba(240, 240, 255, 0.85));
            backdrop-filter: blur(10px);
            border-radius: 30px;
            padding: 50px;
            box-shadow: 0 20px 50px rgba(0, 0, 0, 0.15), inset 0 0 30px rgba(106, 17, 203, 0.1);
            max-width: 90%;
            width: 650px;
            position: relative;
            z-index: 1;
            animation: float 6s infinite ease-in-out;
            border: 1px solid rgba(255, 255, 255, 0.5);
            overflow: hidden;
        }
        
        @keyframes float {
            0%, 100% { transform: translateY(0); }
            50% { transform: translateY(-20px); }
        }
        
        .error-code {
            font-size: 180px;
            font-weight: 900;
            margin-bottom: 0;
            background: linear-gradient(45deg, #6a11cb, #2575fc, #6a11cb);
            background-size: 200% 200%;
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            line-height: 1;
            text-shadow: 0 5px 15px rgba(106, 17, 203, 0.2);
            animation: gradient 8s ease infinite;
            letter-spacing: -5px;
            position: relative;
        }
        
        @keyframes gradient {
            0% { background-position: 0% 50%; }
            50% { background-position: 100% 50%; }
            100% { background-position: 0% 50%; }
        }
        
        .error-text {
            font-size: 32px;
            font-weight: 700;
            margin: 10px 0 30px;
            color: #444;
            position: relative;
            display: inline-block;
            padding: 0 15px;
        }
        
        .error-text::before,
        .error-text::after {
            content: '';
            position: absolute;
            height: 2px;
            width: 30px;
            background: linear-gradient(45deg, #6a11cb, #2575fc);
            top: 50%;
            transform: translateY(-50%);
        }
        
        .error-text::before {
            left: -40px;
        }
        
        .error-text::after {
            right: -40px;
        }
        
        .error-description {
            font-size: 20px;
            color: #555;
            margin-bottom: 20px;
            line-height: 1.8;
            background: linear-gradient(45deg, #555, #888);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            padding: 10px 20px;
            position: relative;
            z-index: 2;
            text-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
        }
        
        .astronaut {
            width: 150px;
            height: 150px;
            margin: 0 auto 30px;
            position: relative;
            animation: astronaut 15s infinite linear;
            filter: drop-shadow(0 10px 20px rgba(106, 17, 203, 0.3));
            transform-origin: center center;
        }
        
        @keyframes astronaut {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
        
        .astronaut img {
            width: 100%;
            height: 100%;
            object-fit: contain;
        }
        
        .cosmic-dust {
            position: absolute;
            width: 100%;
            height: 100%;
            top: 0;
            left: 0;
            pointer-events: none;
        }
        
        .dust-particle {
            position: absolute;
            width: 3px;
            height: 3px;
            background-color: rgba(255, 255, 255, 0.5);
            border-radius: 50%;
            animation: float-dust 15s infinite linear;
        }
        
        @keyframes float-dust {
            0% {
                transform: translateY(0) rotate(0deg);
                opacity: 0;
            }
            10% {
                opacity: 1;
            }
            90% {
                opacity: 1;
            }
            100% {
                transform: translateY(-100px) rotate(360deg);
                opacity: 0;
            }
        }
        
        .planet {
            position: absolute;
            border-radius: 50%;
            opacity: 0.6;
        }
        
        .planet-1 {
            width: 80px;
            height: 80px;
            background: linear-gradient(45deg, #ff9a9e, #fad0c4);
            top: 10%;
            left: 10%;
            animation: rotate 20s infinite linear;
        }
        
        .planet-2 {
            width: 60px;
            height: 60px;
            background: linear-gradient(45deg, #a1c4fd, #c2e9fb);
            bottom: 15%;
            right: 10%;
            animation: rotate 15s infinite linear reverse;
        }
        
        .planet-3 {
            width: 40px;
            height: 40px;
            background: linear-gradient(45deg, #d4fc79, #96e6a1);
            top: 20%;
            right: 15%;
            animation: rotate 25s infinite linear;
        }
        
        @keyframes rotate {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
        
        .meteor {
            position: absolute;
            width: 300px;
            height: 1px;
            transform: rotate(-45deg);
            background: linear-gradient(90deg, rgba(255,255,255,0), rgba(255,255,255,1), rgba(255,255,255,0));
            animation: meteor 5s infinite ease-in-out;
            opacity: 0;
        }
        
        .meteor:nth-child(1) {
            top: 20%;
            left: 80%;
            animation-delay: 0s;
        }
        
        .meteor:nth-child(2) {
            top: 60%;
            left: 90%;
            animation-delay: 2s;
        }
        
        .meteor:nth-child(3) {
            top: 40%;
            left: 70%;
            animation-delay: 4s;
        }
        
        @keyframes meteor {
            0%, 100% { 
                transform: translateX(0) translateY(0) rotate(-45deg);
                opacity: 0;
            }
            10%, 90% { opacity: 1; }
            100% { 
                transform: translateX(-500px) translateY(500px) rotate(-45deg);
            }
        }
        
        @media (max-width: 768px) {
            .error-code {
                font-size: 120px;
            }
            
            .error-text {
                font-size: 24px;
            }
            
            .container {
                padding: 30px;
            }
            
            .astronaut {
                width: 100px;
                height: 100px;
            }
            
            .error-text::before,
            .error-text::after {
                width: 20px;
            }
            
            .error-text::before {
                left: -25px;
            }
            
            .error-text::after {
                right: -25px;
            }
        }
    </style>
</head>
<body>
    <div class="stars" id="stars"></div>
    
    <div class="planet planet-1"></div>
    <div class="planet planet-2"></div>
    <div class="planet planet-3"></div>
    
    <div class="meteor"></div>
    <div class="meteor"></div>
    <div class="meteor"></div>
    
    <div class="container">
        <div class="astronaut">
            <svg viewBox="0 0 512 512" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M256 56C162.904 56 88 130.904 88 224C88 317.096 162.904 392 256 392C349.096 392 424 317.096 424 224C424 130.904 349.096 56 256 56ZM256 376C171.738 376 104 308.262 104 224C104 139.738 171.738 72 256 72C340.262 72 408 139.738 408 224C408 308.262 340.262 376 256 376Z" fill="#6A11CB"/>
                <path d="M256 96C185.307 96 128 153.307 128 224C128 294.693 185.307 352 256 352C326.693 352 384 294.693 384 224C384 153.307 326.693 96 256 96ZM256 336C194.144 336 144 285.856 144 224C144 162.144 194.144 112 256 112C317.856 112 368 162.144 368 224C368 285.856 317.856 336 256 336Z" fill="#2575FC"/>
                <path d="M256 144C210.112 144 176 178.112 176 224C176 269.888 210.112 304 256 304C301.888 304 336 269.888 336 224C336 178.112 301.888 144 256 144ZM256 288C219.944 288 192 260.056 192 224C192 187.944 219.944 160 256 160C292.056 160 320 187.944 320 224C320 260.056 292.056 288 256 288Z" fill="#6A11CB"/>
                <path d="M256 192C236.118 192 220 208.118 220 228C220 247.882 236.118 264 256 264C275.882 264 292 247.882 292 228C292 208.118 275.882 192 256 192ZM256 248C245.028 248 236 238.972 236 228C236 217.028 245.028 208 256 208C266.972 208 276 217.028 276 228C276 238.972 266.972 248 256 248Z" fill="#2575FC"/>
                <path d="M256 0C251.582 0 248 3.582 248 8V40C248 44.418 251.582 48 256 48C260.418 48 264 44.418 264 40V8C264 3.582 260.418 0 256 0Z" fill="#6A11CB"/>
                <path d="M256 400C251.582 400 248 403.582 248 408V504C248 508.418 251.582 512 256 512C260.418 512 264 508.418 264 504V408C264 403.582 260.418 400 256 400Z" fill="#6A11CB"/>
                <path d="M504 248H408C403.582 248 400 251.582 400 256C400 260.418 403.582 264 408 264H504C508.418 264 512 260.418 512 256C512 251.582 508.418 248 504 248Z" fill="#6A11CB"/>
                <path d="M104 248H8C3.582 248 0 251.582 0 256C0 260.418 3.582 264 8 264H104C108.418 264 112 260.418 112 256C112 251.582 108.418 248 104 248Z" fill="#6A11CB"/>
                <path d="M120.686 120.686C117.372 124.001 117.372 129.333 120.686 132.648L162.75 174.711C166.065 178.026 171.397 178.026 174.711 174.711C178.026 171.397 178.026 166.065 174.711 162.75L132.648 120.686C129.333 117.372 124.001 117.372 120.686 120.686Z" fill="#2575FC"/>
                <path d="M337.289 337.289C333.975 340.604 333.975 345.936 337.289 349.251L379.352 391.314C382.667 394.629 387.999 394.629 391.314 391.314C394.628 387.999 394.628 382.667 391.314 379.352L349.251 337.289C345.936 333.975 340.604 333.975 337.289 337.289Z" fill="#2575FC"/>
                <path d="M391.314 120.686C387.999 117.372 382.667 117.372 379.352 120.686L337.289 162.75C333.975 166.065 333.975 171.397 337.289 174.711C340.604 178.026 345.936 178.026 349.251 174.711L391.314 132.648C394.628 129.333 394.628 124.001 391.314 120.686Z" fill="#2575FC"/>
                <path d="M174.711 337.289C171.397 333.975 166.065 333.975 162.75 337.289L120.686 379.352C117.372 382.667 117.372 387.999 120.686 391.314C124.001 394.629 129.333 394.629 132.648 391.314L174.711 349.251C178.026 345.936 178.026 340.604 174.711 337.289Z" fill="#2575FC"/>
            </svg>
        </div>
        <h1 class="error-code">404</h1>
        <h2 class="error-text">页面未找到</h2>
        <p class="error-description">哎呀！看起来您要寻找的页面已经飞向了太空。<br>我们的宇航员正在努力寻找它，但似乎已经迷失在浩瀚的宇宙中。</p>
        <div class="cosmic-dust" id="cosmicDust"></div>
    </div>

    <script>
        // 创建星星背景
        function createStars() {
            const stars = document.getElementById('stars');
            const count = 300;
            
            for (let i = 0; i < count; i++) {
                const star = document.createElement('div');
                star.className = 'star';
                
                // 随机位置
                const x = Math.random() * 100;
                const y = Math.random() * 100;
                star.style.left = x + '%';
                star.style.top = y + '%';
                
                // 随机大小
                const size = Math.random() * 4;
                star.style.width = size + 'px';
                star.style.height = size + 'px';
                
                // 随机动画持续时间和不透明度
                const duration = 3 + Math.random() * 7;
                const opacity = 0.2 + Math.random() * 0.8;
                star.style.setProperty('--duration', duration + 's');
                star.style.setProperty('--opacity', opacity);
                
                stars.appendChild(star);
            }
        }
        
        // 创建宇宙尘埃
        function createCosmicDust() {
            const container = document.getElementById('cosmicDust');
            const count = 30;
            
            for (let i = 0; i < count; i++) {
                const dust = document.createElement('div');
                dust.className = 'dust-particle';
                
                // 随机位置
                const x = Math.random() * 100;
                const y = Math.random() * 100;
                dust.style.left = x + '%';
                dust.style.top = y + '%';
                
                // 随机大小
                const size = Math.random() * 2 + 1;
                dust.style.width = size + 'px';
                dust.style.height = size + 'px';
                
                // 随机动画延迟
                const delay = Math.random() * 10;
                dust.style.animationDelay = delay + 's';
                
                container.appendChild(dust);
            }
        }
        
        // 页面加载时创建星星和宇宙尘埃
        window.addEventListener('load', function() {
            createStars();
            createCosmicDust();
        });
    </script>
</body>
</html>`

	// Html500 StatusInternalServerError
	Html500 = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>500 - 系统错误</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Microsoft YaHei', Arial, sans-serif;
            background-color: #f0f4f8;
            color: #333;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            overflow: hidden;
            position: relative;
        }
        
        .background {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            z-index: -1;
            overflow: hidden;
        }
        
        .bubble {
            position: absolute;
            border-radius: 50%;
            background: linear-gradient(45deg, rgba(255,255,255,0.1), rgba(255,255,255,0.3));
            animation: float var(--duration) infinite ease-in-out;
            opacity: 0.6;
        }
        
        @keyframes float {
            0%, 100% { transform: translateY(0) translateX(0); }
            25% { transform: translateY(-20px) translateX(10px); }
            50% { transform: translateY(0) translateX(20px); }
            75% { transform: translateY(20px) translateX(10px); }
        }
        
        .container {
            text-align: center;
            background: linear-gradient(135deg, rgba(255, 255, 255, 0.9), rgba(240, 248, 255, 0.85));
            backdrop-filter: blur(10px);
            border-radius: 30px;
            padding: 50px;
            box-shadow: 0 20px 50px rgba(0, 0, 0, 0.1), inset 0 0 30px rgba(66, 153, 225, 0.1);
            max-width: 90%;
            width: 650px;
            position: relative;
            z-index: 1;
            animation: pulse 6s infinite ease-in-out;
            border: 1px solid rgba(255, 255, 255, 0.5);
            overflow: hidden;
        }
        
        @keyframes pulse {
            0%, 100% { transform: scale(1); }
            50% { transform: scale(1.02); }
        }
        
        .error-code {
            font-size: 180px;
            font-weight: 900;
            margin-bottom: 0;
            background: linear-gradient(45deg, #3182ce, #63b3ed, #3182ce);
            background-size: 200% 200%;
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            line-height: 1;
            text-shadow: 0 5px 15px rgba(49, 130, 206, 0.2);
            animation: gradient 8s ease infinite;
            letter-spacing: -5px;
            position: relative;
        }
        
        @keyframes gradient {
            0% { background-position: 0% 50%; }
            50% { background-position: 100% 50%; }
            100% { background-position: 0% 50%; }
        }
        
        .error-text {
            font-size: 32px;
            font-weight: 700;
            margin: 10px 0 30px;
            color: #2c5282;
            position: relative;
            display: inline-block;
            padding: 0 15px;
        }
        
        .error-text::before,
        .error-text::after {
            content: '';
            position: absolute;
            height: 2px;
            width: 30px;
            background: linear-gradient(45deg, #3182ce, #63b3ed);
            top: 50%;
            transform: translateY(-50%);
        }
        
        .error-text::before {
            left: -40px;
        }
        
        .error-text::after {
            right: -40px;
        }
        
        .error-description {
            font-size: 20px;
            color: #4a5568;
            margin-bottom: 20px;
            line-height: 1.8;
            background: linear-gradient(45deg, #4a5568, #718096);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            padding: 10px 20px;
            position: relative;
            z-index: 2;
            text-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
        }
        
        .robot {
            width: 150px;
            height: 150px;
            margin: 0 auto 30px;
            position: relative;
            filter: drop-shadow(0 10px 20px rgba(49, 130, 206, 0.3));
        }
        
        .robot svg {
            width: 100%;
            height: 100%;
            animation: wobble 3s infinite ease-in-out;
            transform-origin: center bottom;
        }
        
        @keyframes wobble {
            0%, 100% { transform: rotate(0deg); }
            25% { transform: rotate(-5deg); }
            75% { transform: rotate(5deg); }
        }
        
        .gear {
            position: absolute;
            border-radius: 50%;
            background: linear-gradient(45deg, #3182ce, #63b3ed);
            display: flex;
            justify-content: center;
            align-items: center;
            animation: spin 10s infinite linear;
        }
        
        .gear::before {
            content: '';
            width: 40%;
            height: 40%;
            background-color: white;
            border-radius: 50%;
        }
        
        .gear-1 {
            width: 60px;
            height: 60px;
            top: 10%;
            left: 10%;
        }
        
        .gear-2 {
            width: 40px;
            height: 40px;
            bottom: 15%;
            right: 10%;
            animation-direction: reverse;
        }
        
        .gear-3 {
            width: 30px;
            height: 30px;
            top: 20%;
            right: 15%;
        }
        
        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
        
        .particles {
            position: absolute;
            width: 100%;
            height: 100%;
            top: 0;
            left: 0;
            pointer-events: none;
        }
        
        .particle {
            position: absolute;
            width: 8px;
            height: 8px;
            background-color: rgba(99, 179, 237, 0.5);
            border-radius: 50%;
            animation: float-particle 15s infinite linear;
        }
        
        @keyframes float-particle {
            0% {
                transform: translateY(0) rotate(0deg);
                opacity: 0;
            }
            10% {
                opacity: 1;
            }
            90% {
                opacity: 1;
            }
            100% {
                transform: translateY(-100px) rotate(360deg);
                opacity: 0;
            }
        }
        
        @media (max-width: 768px) {
            .error-code {
                font-size: 120px;
            }
            
            .error-text {
                font-size: 24px;
            }
            
            .container {
                padding: 30px;
            }
            
            .robot {
                width: 100px;
                height: 100px;
            }
            
            .error-text::before,
            .error-text::after {
                width: 20px;
            }
            
            .error-text::before {
                left: -25px;
            }
            
            .error-text::after {
                right: -25px;
            }
            
            .gear-1 {
                width: 40px;
                height: 40px;
            }
            
            .gear-2 {
                width: 30px;
                height: 30px;
            }
            
            .gear-3 {
                width: 20px;
                height: 20px;
            }
        }
    </style>
</head>
<body>
    <div class="background" id="background"></div>
    
    <div class="gear gear-1"></div>
    <div class="gear gear-2"></div>
    <div class="gear gear-3"></div>
    
    <div class="container">
        <div class="robot">
            <svg viewBox="0 0 512 512" fill="none" xmlns="http://www.w3.org/2000/svg">
                <rect x="128" y="128" width="256" height="256" rx="16" fill="#3182CE"/>
                <rect x="168" y="168" width="176" height="176" rx="8" fill="white"/>
                <circle cx="208" cy="216" r="24" fill="#3182CE"/>
                <circle cx="304" cy="216" r="24" fill="#3182CE"/>
                <path d="M208 296H304" stroke="#3182CE" stroke-width="16" stroke-linecap="round"/>
                <path d="M96 208H128" stroke="#3182CE" stroke-width="16" stroke-linecap="round"/>
                <path d="M384 208H416" stroke="#3182CE" stroke-width="16" stroke-linecap="round"/>
                <path d="M96 304H128" stroke="#3182CE" stroke-width="16" stroke-linecap="round"/>
                <path d="M384 304H416" stroke="#3182CE" stroke-width="16" stroke-linecap="round"/>
                <path d="M208 384V416" stroke="#3182CE" stroke-width="16" stroke-linecap="round"/>
                <path d="M256 384V416" stroke="#3182CE" stroke-width="16" stroke-linecap="round"/>
                <path d="M304 384V416" stroke="#3182CE" stroke-width="16" stroke-linecap="round"/>
                <path d="M208 96V128" stroke="#3182CE" stroke-width="16" stroke-linecap="round"/>
                <path d="M256 96V128" stroke="#3182CE" stroke-width="16" stroke-linecap="round"/>
                <path d="M304 96V128" stroke="#3182CE" stroke-width="16" stroke-linecap="round"/>
            </svg>
        </div>
        <h1 class="error-code">500</h1>
        <h2 class="error-text">系统错误</h2>
        <p class="error-description">哎呀！系统开小差啦，请稍后再试。<br>我们的工程师正在努力修复这个问题。</p>
        <div class="particles" id="particles"></div>
    </div>

    <script>
        // 创建气泡背景
        function createBubbles() {
            const background = document.getElementById('background');
            const count = 20;
            
            for (let i = 0; i < count; i++) {
                const bubble = document.createElement('div');
                bubble.className = 'bubble';
                
                // 随机位置
                const x = Math.random() * 100;
                const y = Math.random() * 100;
                bubble.style.left = x + '%';
                bubble.style.top = y + '%';
                
                // 随机大小
                const size = 50 + Math.random() * 100;
                bubble.style.width = size + 'px';
                bubble.style.height = size + 'px';
                
                // 随机动画持续时间
                const duration = 10 + Math.random() * 20;
                bubble.style.setProperty('--duration', duration + 's');
                
                // 随机延迟
                const delay = Math.random() * 10;
                bubble.style.animationDelay = delay + 's';
                
                background.appendChild(bubble);
            }
        }
        
        // 创建粒子
        function createParticles() {
            const container = document.getElementById('particles');
            const count = 20;
            
            for (let i = 0; i < count; i++) {
                const particle = document.createElement('div');
                particle.className = 'particle';
                
                // 随机位置
                const x = Math.random() * 100;
                const y = Math.random() * 100;
                particle.style.left = x + '%';
                particle.style.top = y + '%';
                
                // 随机大小
                const size = Math.random() * 5 + 3;
                particle.style.width = size + 'px';
                particle.style.height = size + 'px';
                
                // 随机动画延迟
                const delay = Math.random() * 10;
                particle.style.animationDelay = delay + 's';
                
                container.appendChild(particle);
            }
        }
        
        // 页面加载时创建元素
        window.addEventListener('load', function() {
            createBubbles();
            createParticles();
        });
    </script>
</body>
</html>`
)
