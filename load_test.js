import http from 'k6/http';
import ws from 'k6/ws';
import { check, sleep, fail } from 'k6';

export const options = {
    scenarios: {
        load_test: {
            executor: 'ramping-vus',
            stages: [
                { duration: '1m', target: 10 }, // Stay under the 30-connection limit
                { duration: '1m', target: 15 },
                { duration: '30s', target: 0 },
            ],
        },
    },
    thresholds: {
        http_req_duration: ['p(95)<1000'], // Increased slightly to account for DB latency
        http_req_failed: ['rate<0.1'],
    },
};

const BASE_URL = 'http://localhost:8081';

// --- SETUP PHASE ---
export function setup() {
    const userName = `SetupUser_${Math.floor(Math.random() * 10000)}`;

    // 1. Create User
    const createRes = http.post(`${BASE_URL}/api/v1/create-user`,
        { user_name: userName },
        { headers: { 'Content-Type': 'application/x-www-form-urlencoded' } }
    );

    if (createRes.status !== 201) {
        fail(`Setup failed: CreateUser returned ${createRes.status}`);
    }

    // 2. Auth User (to get cookie)
    const authRes = http.post(`${BASE_URL}/api/v1/auth-user`,
        JSON.stringify({ Name: userName }),
        { headers: { 'Content-Type': 'application/json' } }
    );

    const sessionId = authRes.cookies.session_id ? authRes.cookies.session_id[0].value : null;
    if (!sessionId) fail('Setup failed: No session cookie received');

    // 3. Create 5 Rooms for testing
    const roomIds = [];
    for (let i = 0; i < 5; i++) {
        const roomRes = http.post(`${BASE_URL}/api/v1/rooms/`,
            { name: `LoadRoom_${i}`, description: 'Test', topic: 'Testing', is_private: 'false', max_members: '50' },
            {
                headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
                cookies: { session_id: sessionId }
            }
        );

        if (roomRes.status === 201) {
            // Logic check: Ensure your Go API returns the ID in the "data" or root object
            const body = JSON.parse(roomRes.body);
            const id = body.data ? body.data.id : body.id;
            if (id) roomIds.push(id);
        }
    }

    if (roomIds.length === 0) fail('Setup failed: No rooms could be created');
    return { roomIds };
}

// --- VIRTUAL USER LOGIC ---
export default function (data) {
    const roomIds = data.roomIds;
    const username = `VU_${__VU}_${Math.floor(Math.random() * 1000)}`;

    // 1. Create & Auth VU
    http.post(`${BASE_URL}/api/v1/create-user`, { user_name: username }, { headers: { 'Content-Type': 'application/x-www-form-urlencoded' } });

    const authRes = http.post(`${BASE_URL}/api/v1/auth-user`,
        JSON.stringify({ Name: username }),
        { headers: { 'Content-Type': 'application/json' } }
    );

    const sessionId = authRes.cookies.session_id ? authRes.cookies.session_id[0].value : null;
    if (!sessionId) return;

    const cookieHeader = { session_id: sessionId };

    // 2. Select a Random Room
    const roomId = roomIds[Math.floor(Math.random() * roomIds.length)];

    // 3. Join Room
    const joinRes = http.post(`${BASE_URL}/api/v1/rooms/join/${roomId}`, {}, { cookies: cookieHeader });
    check(joinRes, { 'joined room': (r) => r.status === 200 });

    // 4. WebSocket Interaction
    const wsUrl = `ws://localhost:8081/api/v1/rooms/ws/${roomId}`;
    const wsRes = ws.connect(wsUrl, { cookies: cookieHeader }, function (socket) {
        socket.on('open', () => {
            socket.send(JSON.stringify({ type: 'message.send', content: `Hello from VU ${__VU}` }));

            // Simulate chatty user: Send a message every 2 seconds
            socket.setInterval(() => {
                socket.send(JSON.stringify({ type: 'message.send', content: 'Still here!' }));
            }, 2000);
        });

        socket.on('error', (e) => console.log(`VU ${__VU} WS Error: ${e.error()}`));

        // Hold connection for 10 seconds to simulate a real user session
        socket.setTimeout(() => socket.close(), 10000);
    });

    check(wsRes, { 'WS connected': (r) => r && r.status === 101 });

    // 5. Cleanup
    sleep(1);
}