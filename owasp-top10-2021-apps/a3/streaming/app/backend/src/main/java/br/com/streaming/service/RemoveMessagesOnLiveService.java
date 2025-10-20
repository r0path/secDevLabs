package br.com.streaming.service;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import br.com.streaming.model.Live;
import br.com.streaming.repository.LiveRepository;

@Service
public class RemoveMessagesOnLiveService {
	
	@Autowired
	private LiveRepository liveRepository;

	public void remove(Long id) {
		liveRepository.findById(id).ifPresent(live -> {
			live.removeMessages();
			liveRepository.save(live);
		});
	}
}
