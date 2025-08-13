package frame

const (
	PageNotFind = `<!DOCTYPE html>
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
)
