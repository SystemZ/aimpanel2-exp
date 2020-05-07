export interface GameServer {
    id: string,
    name: string,
    host_id: string,
    state: number,
    state_last_changed: Date,
    game_id: number,
    game_version: string,
    metric_frequency: number,
    stop_timeout: number
    created_at: Date,
}

export interface GameServerList {
    game_servers: GameServer[]
}
