using System.Text.Json;

using Wallpaper.CommonLanguage;

namespace WallpaperUI.Cs
{
    public class PreferencesManager2
    {
        private List<Func<Dictionary<string, dynamic>,Dictionary<string, dynamic>>> changeHandler = new();
        private Dictionary<string, dynamic> _preferences;
        public bool loopBreak = false;

        public PreferencesManager2()
        {
            _preferences = new Dictionary<string, dynamic>();

        }
        public PreferencesManager2(Dictionary<string, dynamic> preferences)
        {
            _preferences = preferences;
        }

        public void AddChangeHandler(Func<Dictionary<string, dynamic>, Dictionary<string, dynamic>> handler)
        {
            changeHandler.Add(handler);
        }

        public void RemoveChangeHandler(Func<Dictionary<string, dynamic>, Dictionary<string, dynamic>> handler)
        {
            changeHandler.Remove(handler);
        }

        public Dictionary<string, dynamic> GetPreferences()
        {
            return _preferences;
        }

        public PreferencesManager2(List<Func<Dictionary<string, dynamic>, Dictionary<string, dynamic>>> changeHandler, Dictionary<string, dynamic> preferences)
        {
            this.changeHandler = changeHandler;
            _preferences = preferences;
        }

        public PreferencesManager2(Dictionary<string, dynamic> preferences, params Func<Dictionary<string, dynamic>, Dictionary<string, dynamic>>[] handlers)
        {
            this.changeHandler = handlers.ToList();
            _preferences = preferences;
        }



        public string Serialized()
        {
            return JsonSerializer.Serialize(_preferences);
        }

        public T Get<T>(string key)
        {
            return _preferences[key];
        }

        public T GetOrDefault<T>(string key, T defaultValue)
        {
            return _preferences.ContainsKey(key) ? _preferences[key] : defaultValue;
        }

        public dynamic GetOrDefault2(string key, dynamic defaultValue)
        {
            return _preferences.ContainsKey(key) ? _preferences[key] : defaultValue;
        }

        public void Set(string key, dynamic value)
        {
            _preferences[key] = value;
            foreach (var handler in changeHandler.ToList())
            {
                _preferences = handler(_preferences);
            }
        }

        public void Remove(string key)
        {
            _preferences.Remove(key);
            foreach (var handler in changeHandler.ToList())
            {
                _preferences = handler(_preferences);
            }
        }

        public void Clear()
        {
            _preferences.Clear();
            foreach (var handler in changeHandler.ToList())
            {
                _preferences = handler(_preferences);
            }
        }

        public void Add(string key, dynamic value)
        {
            if (!_preferences.ContainsKey(key))
            {
                _preferences.Add(key, value);
            }
            foreach (var handler in changeHandler.ToList())
            {
                _preferences = handler(_preferences);
            }
        }

        public void SetFrom(Dictionary<string, dynamic> preferences)
        {
            _preferences = preferences;
            foreach (var handler in changeHandler.ToList())
            {
                _preferences = handler(_preferences);
            }
        }
    }

    public class PreferencesManager
    {
        
        private Wallpaper.CommonLanguage.AppSettings appSettings;
        private List<Func<AppSettings, Task<AppSettings>>> writeHandlers;

        public AppSettings Value => applyDefaults(appSettings.Clone());

        public SimpleBackgroundsSystem SimpleBackgroundsSystem => Value.SimpleBackgroundsSystem == null ? new SimpleBackgroundsSystem() : Value.SimpleBackgroundsSystem;

        public RecentBackgroundStore RecentBackgroundStore => Value.RecentBackgroundStore == null ? new RecentBackgroundStore() : Value.RecentBackgroundStore;

        public RecentColorSystem RecentColorSystem => Value.RecentColorSystem == null ? new RecentColorSystem() : Value.RecentColorSystem;

        public List<RuntimePanel> PanelsToSkip => Value.PanelsToSkip == null ? [] : Value.PanelsToSkip.ToList();

        private AppSettings applyDefaults(AppSettings settings)
        {
            if (settings.RecentBackgroundStore == null)
            {
                settings.RecentBackgroundStore = new RecentBackgroundStore();
            }
            if (settings.RecentColorSystem == null)
            {
                settings.RecentColorSystem = new RecentColorSystem();
            }
            if (settings.SimpleBackgroundsSystem == null)
            {
                settings.SimpleBackgroundsSystem = new SimpleBackgroundsSystem();
            }
            return settings;
        }

        public async Task Write(Func<AppSettings,Task<AppSettings>> func)
        {
            var tmp = applyDefaults(appSettings.Clone());
            appSettings = await func(tmp);

            foreach (var item in writeHandlers)
            {
                appSettings = await item(appSettings.Clone());
            }
        }

        public PreferencesManager()
        {
            appSettings = new AppSettings();
            appSettings.RecentBackgroundStore = new RecentBackgroundStore();
            appSettings.RecentColorSystem  = new RecentColorSystem();
            appSettings.SimpleBackgroundsSystem = new SimpleBackgroundsSystem();
            writeHandlers = new();
        }

        public Action AddWriteHandler(Func<AppSettings, Task<AppSettings>> func)
        {
            writeHandlers.Add(func);
            return () => writeHandlers.Remove(func);
        }
    }
}
