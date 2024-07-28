using System.Text.Json;

namespace WallpaperUI.Cs
{
    public class PreferencesManager
    {
        private List<Func<Dictionary<string, dynamic>,Dictionary<string, dynamic>>> changeHandler = new();
        private Dictionary<string, dynamic> _preferences;
        public bool loopBreak = false;

        public PreferencesManager()
        {
            _preferences = new Dictionary<string, dynamic>();

        }
        public PreferencesManager(Dictionary<string, dynamic> preferences)
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

        public PreferencesManager(List<Func<Dictionary<string, dynamic>, Dictionary<string, dynamic>>> changeHandler, Dictionary<string, dynamic> preferences)
        {
            this.changeHandler = changeHandler;
            _preferences = preferences;
        }

        public PreferencesManager(Dictionary<string, dynamic> preferences, params Func<Dictionary<string, dynamic>, Dictionary<string, dynamic>>[] handlers)
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
}
