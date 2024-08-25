// Prevents additional console window on Windows in release, DO NOT REMOVE!!
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

#[tauri::command]
fn add_mod(mod_path: &str) -> String {
    println!("add mod {}", mod_path);
    format!("add mod, {}", mod_path)
}

fn main() {
    tauri::Builder::default()
        .invoke_handler(tauri::generate_handler![add_mod])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
