export default function VideoType(type) {
    const v_map = {
        'knowledge': 1,
        'hotpot': 2,
        'game': 3,
        'entertainment': 4,
        'fantasy': 5, 
        'music': 6,
        'food': 7,
        'sport': 8,
        'fashion': 9
    }
    
    return v_map[type] ? v_map[type] : 999;
  }
  